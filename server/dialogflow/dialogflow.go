package dialogflow

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/niko-cb/covid19datascraper/server/env"

	"github.com/niko-cb/covid19datascraper/server/datastore"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"github.com/niko-cb/covid19datascraper/server/model"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type Processor struct {
	projectID      string
	authentication string
	language       string
	timezone       string
	sessionClient  *dialogflow.SessionsClient
	ctx            context.Context
}

var p Processor

func (p *Processor) init() (err error) {
	p.projectID = env.ProjectID()
	p.authentication = env.AuthDialogflow()
	p.language = env.Language()
	p.timezone = env.Timezone()

	p.ctx = context.Background()
	sessionClient, err := dialogflow.NewSessionsClient(p.ctx, option.WithCredentialsJSON([]byte(p.authentication)))
	if err != nil {
		log.Fatalf("Failed to authenticate with Dialogflow: %v", err)
	}

	p.sessionClient = sessionClient
	return
}

func (p *Processor) CreateOrRecreateIntents() error {
	ctx := p.ctx
	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsJSON([]byte(p.authentication)))
	if clientErr != nil {
		return clientErr
	}
	defer intentsClient.Close()

	parent := fmt.Sprintf("projects/%s/agent", p.projectID)

	if err := deleteIntents(ctx, intentsClient, p.projectID, parent); err != nil {
		return err
	}

	if err := addPrefectureDataIntents(ctx, intentsClient, parent); err != nil {
		return err
	}

	if err := addCoronavirusSymptomsIntent(ctx, intentsClient, parent); err != nil {
		return err
	}

	return nil
}

func addPrefectureDataIntents(ctx context.Context, intentsClient *dialogflow.IntentsClient, parent string) error {
	dsClient, err := datastore.NewClient()
	if err != nil {
		return err
	}
	pData, err := model.GetPrefectureData(ctx, dsClient, datastore.DataKind())
	if err != nil {
		return err
	}
	sd, err := model.GetDateFromDatastore(ctx, dsClient, datastore.DateKind())
	if err != nil {
		return err
	}
	for _, p := range pData {
		// 1 second sleep in order to prevent hitting the request limit
		time.Sleep(1 * time.Second)
		displayName := p.Prefecture
		var trainingPhraseParts []string
		trainingPhraseParts = append(trainingPhraseParts, displayName)
		var messageTexts []string
		cities := "\n\n" + "[市町村の公開データ]\n"
		cityData := cityMap(p.ConfirmedByCity)
		if cityData != "" {
			cities = cities + cityData
		}
		messageTexts = append(messageTexts, sd.Date+"までの情報です\n\n"+"都道府県名:   "+displayName+"\n検査陽性者:   "+p.Cases+"\n回復者:   "+p.Recovered+"\n死者:   "+p.Deaths+"\n昨日の件数:  "+p.YesterdayConfirmed+"\n今日の件数:  "+p.NewlyConfirmed+cities)

		request := createDialogflowIntent(displayName, parent, trainingPhraseParts, messageTexts)

		_, requestErr := intentsClient.CreateIntent(ctx, &request)
		log.Println(&request)
		if requestErr != nil {
			return requestErr
		}
	}
	return nil
}

func addCoronavirusSymptomsIntent(ctx context.Context, intentsClient *dialogflow.IntentsClient, parent string) error {
	cs := model.GetCoronavirusSymptoms()

	commonString := strings.Join(cs.Common, ", ")
	rareString := strings.Join(cs.Rare, ", ")
	severeString := strings.Join(cs.Severe, ", ")

	displayName := "症状"
	var trainingPhraseParts []string
	trainingPhraseParts = append(trainingPhraseParts, displayName)
	var messageTexts []string
	messageTexts = append(messageTexts, "症状\n\n"+"初期症状: "+commonString+"\n\n"+"人によっての症状: "+rareString+"\n\n"+"重篤な症状: "+severeString)
	request := createDialogflowIntent(displayName, parent, trainingPhraseParts, messageTexts)

	_, requestErr := intentsClient.CreateIntent(ctx, &request)
	log.Println(&request)
	if requestErr != nil {
		return requestErr
	}
	return nil
}

func listIntents(ctx context.Context, intentsClient *dialogflow.IntentsClient, parent string) ([]*dialogflowpb.Intent, error) {

	request := dialogflowpb.ListIntentsRequest{Parent: parent}

	intentIterator := intentsClient.ListIntents(ctx, &request)
	var intents []*dialogflowpb.Intent

	for intent, status := intentIterator.Next(); status != iterator.Done; {
		intents = append(intents, intent)
		intent, status = intentIterator.Next()
	}

	return intents, nil
}

func deleteIntents(ctx context.Context, intentsClient *dialogflow.IntentsClient, projectID, parent string) error {
	dfIntents, err := listIntents(ctx, intentsClient, parent)
	if err != nil {
		return err
	}
	for _, intent := range dfIntents {
		time.Sleep(1 * time.Second)
		route := strings.Split(intent.GetName(), "/")
		intentID := route[len(route)-1]
		targetPath := fmt.Sprintf("projects/%s/agent/intents/%s", projectID, intentID)
		request := dialogflowpb.DeleteIntentRequest{Name: targetPath}

		requestErr := intentsClient.DeleteIntent(ctx, &request)
		if requestErr != nil {
			return requestErr
		}
	}
	return nil
}

func createDialogflowIntent(displayName, parent string, trainingPhraseParts, messageTexts []string) dialogflowpb.CreateIntentRequest {
	var targetTrainingPhrases []*dialogflowpb.Intent_TrainingPhrase
	var targetTrainingPhraseParts []*dialogflowpb.Intent_TrainingPhrase_Part
	for _, partString := range trainingPhraseParts {
		part := dialogflowpb.Intent_TrainingPhrase_Part{Text: partString}
		targetTrainingPhraseParts = []*dialogflowpb.Intent_TrainingPhrase_Part{&part}
		targetTrainingPhrase := dialogflowpb.Intent_TrainingPhrase{Type: dialogflowpb.Intent_TrainingPhrase_TYPE_UNSPECIFIED, Parts: targetTrainingPhraseParts}
		targetTrainingPhrases = append(targetTrainingPhrases, &targetTrainingPhrase)
	}

	intentMessageTexts := dialogflowpb.Intent_Message_Text{Text: messageTexts}
	wrappedIntentMessageTexts := dialogflowpb.Intent_Message_Text_{Text: &intentMessageTexts}
	intentMessage := dialogflowpb.Intent_Message{Message: &wrappedIntentMessageTexts}

	target := dialogflowpb.Intent{DisplayName: displayName, WebhookState: dialogflowpb.Intent_WEBHOOK_STATE_UNSPECIFIED, TrainingPhrases: targetTrainingPhrases, Messages: []*dialogflowpb.Intent_Message{&intentMessage}}

	return dialogflowpb.CreateIntentRequest{Parent: parent, Intent: &target}
}

func cityMap(cityData string) string {
	var cities []string
	cityList := strings.Split(cityData, ",\"")
	for _, city := range cityList {
		c := strings.ReplaceAll(city, "{", "")
		c = strings.ReplaceAll(c, "}", "")
		c = strings.ReplaceAll(c, "\"", "")
		c = strings.TrimSpace(c)
		c = strings.ReplaceAll(c, ":", ": ")
		if c != "" {
			cities = append(cities, c)
		}
	}
	return strings.Join(cities, "\n")
}

func NewSession() Processor {
	err := p.init()
	if err != nil {
		panic(err)
	}
	return p
}
