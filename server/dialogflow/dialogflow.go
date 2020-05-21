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

func (p *Processor) NewSession() *Processor {
	config := env.Get()
	sessionClient, err := dialogflow.NewSessionsClient(p.ctx, option.WithCredentialsJSON([]byte(config.DialogflowAuth)))
	if err != nil {
		log.Fatalf("Failed to authenticate with Dialogflow: %v", err)
	}
	return &Processor{
		projectID:      config.ProjectID,
		authentication: config.DialogflowAuth,
		language:       config.Language,
		timezone:       config.Timezone,
		sessionClient:  sessionClient,
	}
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
		intent := (&Intent{}).Make(parent, nil, p, sd)
		request := createIntentRequest(intent)
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
	intent := (&Intent{}).Make(parent, cs, nil, nil)
	request := createIntentRequest(intent)

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

func createIntentRequest(intent *Intent) dialogflowpb.CreateIntentRequest {
	var targetTrainingPhrases []*dialogflowpb.Intent_TrainingPhrase
	var targetTrainingPhraseParts []*dialogflowpb.Intent_TrainingPhrase_Part
	for _, partString := range intent.TrainingPhrases {
		part := dialogflowpb.Intent_TrainingPhrase_Part{Text: partString}
		targetTrainingPhraseParts = []*dialogflowpb.Intent_TrainingPhrase_Part{&part}
		targetTrainingPhrase := dialogflowpb.Intent_TrainingPhrase{Type: dialogflowpb.Intent_TrainingPhrase_TYPE_UNSPECIFIED, Parts: targetTrainingPhraseParts}
		targetTrainingPhrases = append(targetTrainingPhrases, &targetTrainingPhrase)
	}

	intentMessageTexts := dialogflowpb.Intent_Message_Text{Text: intent.Messages}
	wrappedIntentMessageTexts := dialogflowpb.Intent_Message_Text_{Text: &intentMessageTexts}
	intentMessage := dialogflowpb.Intent_Message{Message: &wrappedIntentMessageTexts}

	target := dialogflowpb.Intent{DisplayName: intent.DisplayName, WebhookState: dialogflowpb.Intent_WEBHOOK_STATE_UNSPECIFIED, TrainingPhrases: targetTrainingPhrases, Messages: []*dialogflowpb.Intent_Message{&intentMessage}}

	return dialogflowpb.CreateIntentRequest{Parent: intent.Parent, Intent: &target}
}
