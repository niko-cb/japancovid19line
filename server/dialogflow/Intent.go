package dialogflow

import (
	"strings"

	"github.com/niko-cb/covid19datascraper/server/model"
)

type Intent struct {
	DisplayName     string
	Parent          string
	Messages        []string
	TrainingPhrases []string
}

func (i *Intent) Make(parent string, cs *model.Symptoms, p *model.PrefectureData, sd *model.SourceDate) *Intent {
	var displayName string
	var messageTexts []string
	var trainingPhraseParts []string

	if cs != nil {
		displayName = "症状"
		trainingPhraseParts = append(trainingPhraseParts, displayName)
		messageTexts = append(messageTexts, NewSMessage(cs).Content.String())
	} else if p != nil {
		displayName = p.Prefecture
		trainingPhraseParts = append(trainingPhraseParts, displayName)
		cities := cityMap(p.ConfirmedByCity)
		messageTexts = append(messageTexts, NewPMessage(p, sd.Date, cities).Content.String())
	}

	return &Intent{
		DisplayName:     displayName,
		Parent:          parent,
		Messages:        messageTexts,
		TrainingPhrases: trainingPhraseParts,
	}
}

func cityMap(cityData string) string {
	var cities []string
	cityList := strings.Split(cityData, ",\"")
	for _, city := range cityList {
		c := trimString(city)
		if c != "" {
			// because Chigesaki and Kamamura don't exist
			if strings.Split(c, ":")[0] != "Chigesaki" && strings.Split(c, ":")[0] != "Kamamura" {
				cities = append(cities, c)
			}
		}
	}
	return strings.Join(cities, "\n")
}

func trimString(s string) string {
	c := strings.ReplaceAll(s, "{", "")
	c = strings.ReplaceAll(c, "}", "")
	c = strings.ReplaceAll(c, "\"", "")
	c = strings.TrimSpace(c)
	c = strings.ReplaceAll(c, ":", ": ")

	return c
}
