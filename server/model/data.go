package model

type AllData struct {
	Prefectures []Data `json:"prefectures"`
}

type Data struct {
	Confirmed           int            `json:"confirmed"`
	CruisePassenger     int            `json:"cruise_passenger"`
	Deaths              int            `json:"deaths"`
	ConfirmedByCity     map[string]int `json:"confirmedByCity"`
	DailyConfirmedCount []int          `json:"dailyConfirmedCount"`
	DailyConfirmedDate  string         `json:"dailyConfirmedDate"`
	NewlyConfirmed      int            `json:"newlyConfirmed"`
	Recovered           int            `json:"recovered"`
	NameJA              string         `json:"name_ja"`
	Name                string         `json:"name"`
}
