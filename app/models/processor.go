package models

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/api/iterator"

	"github.com/niko-cb/covid19datascraper/app/config/env"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	ctx "github.com/niko-cb/covid19datascraper/app/config/context"
	"github.com/niko-cb/covid19datascraper/app/config/datastore"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type Processor struct {
	projectID      string
	authentication string
	language       string
	timezone       string
	sessionClient  *dialogflow.SessionsClient
	intentsClient  *dialogflow.IntentsClient
	parent         string
	ctx            context.Context
}

var p Processor

func NewSession() Processor {
	config := env.Get()

	sessionClient, err := dialogflow.NewSessionsClient(ctx.Get(), option.WithCredentialsJSON([]byte(config.DialogflowAuth)))
	if err != nil {
		log.Fatalf("Failed to authenticate with Dialogflow: %v", err)
	}
	intentsClient, err := dialogflow.NewIntentsClient(ctx.Get(), option.WithCredentialsJSON([]byte(config.DialogflowAuth)))
	if err != nil {
		log.Fatalf("Failed to authenticate with Dialogflow: %v", err)
	}

	p.projectID = config.ProjectID
	p.authentication = config.DialogflowAuth
	p.language = config.Language
	p.timezone = config.Timezone
	p.sessionClient = sessionClient
	p.intentsClient = intentsClient
	p.parent = fmt.Sprintf("projects/%s/agent", p.projectID)
	p.ctx = ctx.Get()

	return p
}

func (p *Processor) CreateDataIntents() error {
	dsClient, err := datastore.NewClient()
	if err != nil {
		return err
	}
	pData, err := GetPrefectureData(p.ctx, dsClient, datastore.DataKind())
	if err != nil {
		return err
	}
	sd, err := GetDateFromDatastore(p.ctx, dsClient, datastore.DateKind())
	if err != nil {
		return err
	}
	for _, data := range pData {
		// 1 second sleep in order to prevent hitting the request limit
		intent := (&Intent{}).Prefecture(p.parent, data, sd)
		request := createIntentRequest(intent)
		_, requestErr := p.intentsClient.CreateIntent(p.ctx, &request)
		log.Println(&request)
		if requestErr != nil {
			return requestErr
		}
	}
	return nil
}

func (p *Processor) CreateSymptomIntents() error {
	cs := (&Symptoms{}).Get()
	intent := (&Intent{}).Symptoms(p.parent, cs)
	request := createIntentRequest(intent)

	_, requestErr := p.intentsClient.CreateIntent(p.ctx, &request)
	log.Println(&request)
	if requestErr != nil {
		return requestErr
	}
	return nil
}

func (p *Processor) DeleteIntents() error {
	dfIntents, err := listIntents(p.ctx, p.intentsClient, p.parent)
	if err != nil {
		return err
	}
	for _, intent := range dfIntents {
		route := strings.Split(intent.GetName(), "/")
		intentID := route[len(route)-1]
		targetPath := fmt.Sprintf("projects/%s/agent/intents/%s", p.projectID, intentID)
		request := dialogflowpb.DeleteIntentRequest{Name: targetPath}

		requestErr := p.intentsClient.DeleteIntent(p.ctx, &request)
		if requestErr != nil {
			return requestErr
		}
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
