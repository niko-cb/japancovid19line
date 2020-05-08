package dialogflow

import "github.com/niko-cb/covid19datascraper/server/model"

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

func NewPMessage(p *model.PrefectureData, date, cities string) *PMessage {
	return &PMessage{
		&PContent{
			Prefecture:         PrefecturePrefix + p.Prefecture,
			Date:               date + DateSuffix,
			Cases:              CasesPrefix + p.Cases,
			NewConfirmed:       TodayConfirmedPrefix + p.NewlyConfirmed,
			YesterdayConfirmed: YesterdayConfirmedPrefix + p.YesterdayConfirmed,
			Recovered:          RecoveredPrefix + p.Recovered,
			Deaths:             DeathsPrefix + p.Deaths,
			Cities:             CitiesPrefix + cities,
		},
	}
}
