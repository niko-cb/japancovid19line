package utils

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"context"
	"fmt"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/niko-cb/covid19datascraper/server/model"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type DialogflowProcessor struct {
	projectID     string
	authJSON      string
	lang          string
	timeZone      string
	sessionClient *dialogflow.SessionsClient
	ctx           context.Context
}

type NLPResponse struct {
	Intent     string            `json:"intent"`
	Confidence float32           `json:"confidence"`
	Entities   map[string]string `json:"entities"`
}

var dp DialogflowProcessor

func (dp *DialogflowProcessor) init() (err error) {
	dp.projectID = os.Getenv("PROJECT_ID")
	dp.authJSON = os.Getenv("DIALOGFLOW_KEYFILE_JSON")
	dp.lang = "ja"
	dp.timeZone = "Japan/Tokyo"

	dp.ctx = context.Background()
	sessionClient, err := dialogflow.NewSessionsClient(dp.ctx, option.WithCredentialsJSON([]byte(dp.authJSON)))
	if err != nil {
		log.Fatalf("Failed to authenticate with Dialogflow: %v", err)
	}

	dp.sessionClient = sessionClient
	return
}

func (dp *DialogflowProcessor) ProcessNLP(rawMessage string, username string) (r NLPResponse) {
	sessionID := username
	request := dialogflowpb.DetectIntentRequest{
		Session: fmt.Sprintf("projects/%s/agent/sessions/%s", dp.projectID, sessionID),
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
					Text:         rawMessage,
					LanguageCode: dp.lang,
				},
			},
		},
		QueryParams: &dialogflowpb.QueryParameters{
			TimeZone: dp.timeZone,
		},
	}
	response, err := dp.sessionClient.DetectIntent(dp.ctx, &request)
	if err != nil {
		log.Fatalf("Error in dialogflow communication: %s", err.Error())
	}
	queryResult := response.GetQueryResult()
	if queryResult.Intent != nil {
		r.Intent = queryResult.Intent.DisplayName
		r.Confidence = queryResult.IntentDetectionConfidence
	}
	r.Entities = make(map[string]string)
	params := queryResult.Parameters.GetFields()
	if len(params) > 0 {
		for param, p := range params {
			fmt.Printf("Param %s: %s (%s)", param, p.GetStringValue(), p.String())
			extractedValue := extractDialogFlowEntities(p)
			r.Entities[param] = extractedValue
		}
	}
	return
}

func extractDialogFlowEntities(p *structpb.Value) (extractedEntity string) {
	kind := p.GetKind()
	switch kind.(type) {
	case *structpb.Value_StringValue:
		return p.GetStringValue()
	case *structpb.Value_NumberValue:
		return strconv.FormatFloat(p.GetNumberValue(), 'f', 6, 64)
	case *structpb.Value_BoolValue:
		return strconv.FormatBool(p.GetBoolValue())
	case *structpb.Value_StructValue:
		s := p.GetStructValue()
		fields := s.GetFields()
		extractedEntity = ""
		for k, v := range fields {
			if k == "amount" {
				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, strconv.FormatFloat(v.GetNumberValue(), 'f', 6, 64))
			}
			if k == "unit" {
				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, v.GetStringValue())
			}
			if k == "date_time" {
				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, v.GetStringValue())
			}
			// TODO: add more entity types
		}
		return extractedEntity
	default:
		return ""
	}
}

func (dp *DialogflowProcessor) CreateOrRecreateIntents() error {
	ctx := dp.ctx
	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsJSON([]byte(dp.authJSON)))
	if clientErr != nil {
		return clientErr
	}
	defer intentsClient.Close()

	parent := fmt.Sprintf("projects/%s/agent", dp.projectID)

	if err := deleteIntents(ctx, intentsClient, dp.projectID, parent); err != nil {
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
	dsClient, err := NewDSClient()
	if err != nil {
		return err
	}
	pData, err := model.GetPrefectureData(ctx, dsClient, DatastoreKind())
	if err != nil {
		return err
	}
	sd, err := model.GetSourceDate(ctx, dsClient, DatastoreDateKind())
	if err != nil {
		return err
	}
	for _, p := range pData {
		time.Sleep(1 * time.Second)
		displayName := p.Prefecture
		var trainingPhraseParts []string
		trainingPhraseParts = append(trainingPhraseParts, displayName)
		var messageTexts []string
		messageTexts = append(messageTexts, sd.Date+"までの情報です\n\n"+"都道府県名:   "+displayName+"\n検査陽性者:   "+p.Cases+"\n回復者:   "+p.Recovered+"\n死者:   "+p.Deaths)

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

func NewDialogflowSession() (DialogflowProcessor, error) {
	err := dp.init()
	if err != nil {
		log.Fatalf("Failed to initiate dialogflow")
	}
	return dp, nil
}
