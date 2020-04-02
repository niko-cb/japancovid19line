package model

type PrefectureData struct {
	Prefecture   string `json:"prefecture"`
	PCRTests     string `json:"pcr_tests"`
	Hospitalized string `json:"hospitalized"`
	Discharged   string `json:"discharged"`
	Deaths       string `json:"deaths"`
}

type PrefectureDataSlice = []*PrefectureData

func NewPrefectureData(pref, pcr, hos, dis, deaths string) *PrefectureData {
	return &PrefectureData{
		Prefecture:   pref,
		PCRTests:     pcr,
		Hospitalized: hos,
		Discharged:   dis,
		Deaths:       deaths,
	}
}
