package models

import (
	"context"
	"strconv"

	"cloud.google.com/go/datastore"
)

type PrefectureData struct {
	Prefecture         string `json:"prefecture"`
	Cases              string `json:"cases"`
	Recovered          string `json:"recovered"`
	Deaths             string `json:"deaths"`
	NewlyConfirmed     string `json:"newlyConfirmed"`
	YesterdayConfirmed string `json:"yesterdayConfirmed"`

	ConfirmedByCity string `json:"confirmedByCity"`
}

type PrefectureDataSlice = []*PrefectureData

func NewPrefectureData(pref string, cases, rec, deaths, newlyConfirmed, yesterdayConfirmed int, cityMap string) *PrefectureData {
	return &PrefectureData{
		Prefecture:         pref,
		Cases:              strconv.Itoa(cases),
		Recovered:          strconv.Itoa(rec),
		Deaths:             strconv.Itoa(deaths),
		NewlyConfirmed:     strconv.Itoa(newlyConfirmed),
		YesterdayConfirmed: strconv.Itoa(yesterdayConfirmed),
		ConfirmedByCity:    cityMap,
	}
}

func GetPrefectureData(ctx context.Context, dsClient *datastore.Client, kind string) ([]*PrefectureData, error) {
	var pData []*PrefectureData
	q := datastore.NewQuery(kind)
	if _, err := dsClient.GetAll(ctx, q, &pData); err != nil {
		return nil, err
	}
	return pData, nil
}
