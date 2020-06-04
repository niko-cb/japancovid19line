package models

import (
	"strings"
)

type Intent struct {
	DisplayName     string
	Parent          string
	Messages        []string
	TrainingPhrases []string
}

func (i *Intent) Prefecture(parent string, p *PrefectureData, sd *SourceDate) *Intent {
	cities := cityMap(p.ConfirmedByCity)

	i.DisplayName = p.Prefecture
	i.Parent = parent
	i.TrainingPhrases = append(i.TrainingPhrases, i.DisplayName)
	i.Messages = append(i.Messages, NewPMessage(p, sd.Date, cities).Content.String())

	return i
}

func (i *Intent) Symptoms(parent string, cs *Symptoms) *Intent {
	i.DisplayName = "症状"
	i.Parent = parent
	i.TrainingPhrases = append(i.TrainingPhrases, i.DisplayName)
	i.Messages = append(i.Messages, NewSMessage(cs).Content.String())

	return i
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
