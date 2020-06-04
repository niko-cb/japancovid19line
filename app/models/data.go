package models

type CovidDataRes struct {
	Prefectures []PrefectureInfo `json:"prefectures"`
}

type PrefectureInfo struct {
	Confirmed           int            `json:"confirmed"`
	CruisePassenger     int            `json:"cruise_passenger"`
	Deaths              int            `json:"deaths"`
	ConfirmedByCity     map[string]int `json:"confirmedByCity"`
	DailyConfirmedCount []int          `json:"dailyConfirmedCount"`
	DailyConfirmedDate  string         `json:"dailyConfirmedDate"`
	NewlyConfirmed      int            `json:"newlyConfirmed"`
	YesterdayConfirmed  int            `json:"yesterdayConfirmed"`
	Recovered           int            `json:"recovered"`
	NameJA              string         `json:"name_ja"`
	Name                string         `json:"name"`
}
