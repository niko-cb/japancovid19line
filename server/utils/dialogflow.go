package utils

import (
	"cloud.google.com/go/datastore"
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"context"
	"fmt"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/niko-cb/covid19datascraper/server/model"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"log"
	"strconv"
)

type DialogflowProcessor struct {
	projectID        string
	authJSONFilePath string
	lang             string
	timeZone         string
	sessionClient    *dialogflow.SessionsClient
	ctx              context.Context
}

type NLPResponse struct {
	Intent     string            `json:"intent"`
	Confidence float32           `json:"confidence"`
	Entities   map[string]string `json:"entities"`
}

var dp DialogflowProcessor

func (dp *DialogflowProcessor) init(a ...string) (err error) {
	dp.projectID = a[0]
	dp.authJSONFilePath = a[1]
	dp.lang = a[2]
	dp.timeZone = a[3]

	dp.ctx = context.Background()
	sessionClient, err := dialogflow.NewSessionsClient(dp.ctx, option.WithCredentialsFile(dp.authJSONFilePath))
	if err != nil {
		log.Fatal("Failed to authenticate with Dialogflow")
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

func (dp *DialogflowProcessor) CreateIntents() error {
	dsClient, err := NewDSClient()
	if err != nil {
		return err
	}
	kind := DatastoreKind()
	var pData []*model.PrefectureData
	q := datastore.NewQuery(kind)
	if _, err := dsClient.GetAll(context.Background(), q, &pData); err != nil {
		log.Printf("errored on get: %v", err)
	}

	ctx := dp.ctx

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile(dp.authJSONFilePath))
	if clientErr != nil {
		return clientErr
	}
	defer intentsClient.Close()

	projectID := "japancovid19"
	parent := fmt.Sprintf("projects/%s/agent", projectID)

	for _, p := range pData {
		displayName := p.Prefecture
		var trainingPhraseParts []string
		trainingPhraseParts = append(trainingPhraseParts, displayName)
		var messageTexts []string
		messageTexts = append(messageTexts, "都道府県名:   "+displayName+"\nＰＣＲ検査陽性者:   "+p.PCRTests+"\n現在は入院等:   "+p.Hospitalized+"\n退院者:   "+p.Discharged+"\n死亡者:   "+p.Deaths)

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

		request := dialogflowpb.CreateIntentRequest{Parent: parent, Intent: &target}

		_, requestErr := intentsClient.CreateIntent(ctx, &request)
		log.Println(&request)
		if requestErr != nil {
			return requestErr
		}
	}
	return nil
}

func NewDialogflowSession(projectID, authJSONFilePath, lang, timeZone string) (DialogflowProcessor, error) {
	err := dp.init(projectID, authJSONFilePath, lang, timeZone)
	if err != nil {
		log.Fatalf("Failed to initiate dialogflow")
	}
	return dp, nil
}
