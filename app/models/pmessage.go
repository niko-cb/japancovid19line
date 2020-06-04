package models

import "github.com/niko-cb/japancovid19line/app/models/helpers/dialogflow"

// PMessage: the P stands for Prefecture
type PMessage struct {
	Content *PContent
}

// PContent: the P stands for Prefecture
type PContent struct {
	Prefecture         string
	Date               string
	Cases              string
	NewConfirmed       string
	YesterdayConfirmed string
	Recovered          string
	Deaths             string
	Cities             string
}

func (c *PContent) String() string {
	return c.Date + c.Prefecture + c.Cases +
		c.NewConfirmed + c.YesterdayConfirmed +
		c.Recovered + c.Deaths + c.Cities
}

func NewPMessage(p *PrefectureData, date, cities string) *PMessage {
	return &PMessage{
		&PContent{
			Prefecture:         dialogflow.PrefecturePrefix + p.Prefecture,
			Date:               date + dialogflow.DateSuffix,
			Cases:              dialogflow.CasesPrefix + p.Cases,
			NewConfirmed:       dialogflow.TodayConfirmedPrefix + p.NewlyConfirmed,
			YesterdayConfirmed: dialogflow.YesterdayConfirmedPrefix + p.YesterdayConfirmed,
			Recovered:          dialogflow.RecoveredPrefix + p.Recovered,
			Deaths:             dialogflow.DeathsPrefix + p.Deaths,
			Cities:             dialogflow.CitiesPrefix + cities,
		},
	}
}
