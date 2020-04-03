package dialogflow

import (
	"cloud.google.com/go/datastore"
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"context"
	"fmt"
	"github.com/niko-cb/covid19datascraper/server/model"
	"github.com/niko-cb/covid19datascraper/server/utils"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"log"
	"strconv"
)

func CreateIntents(dp utils.DialogflowProcessor) error {
	dsClient, err := utils.NewDSClient()
	if err != nil {
		return err
	}
	kind := utils.DatastoreKind()
	var pData []*model.PrefectureData
	q := datastore.NewQuery(kind)
	if _, err := dsClient.GetAll(context.Background(), q, &pData); err != nil {
		log.Printf("errored on get: %v", err)
	}

	ctx := context.Background()

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx)
	if clientErr != nil {
		return clientErr
	}
	defer intentsClient.Close()

	projectID := "japancovid19"
	parent := fmt.Sprintf("projects/%s/agent", projectID)

	for i, p := range pData {
		displayName := "testing" + strconv.Itoa(i)
		var trainingPhraseParts []string
		trainingPhraseParts = append(trainingPhraseParts, displayName)
		var messageTexts []string
		messageTexts = append(messageTexts, displayName)
		messageTexts = append(messageTexts, p.PCRTests)
		messageTexts = append(messageTexts, p.Hospitalized)
		messageTexts = append(messageTexts, p.Discharged)
		messageTexts = append(messageTexts, p.Deaths)

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

		return nil
	}
	return nil
}
