package model

import (
	"cloud.google.com/go/datastore"
	"context"
)

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

func GetPrefectureData(ctx context.Context, dsClient *datastore.Client, kind string) ([]*PrefectureData, error) {
	var pData []*PrefectureData
	q := datastore.NewQuery(kind)
	if _, err := dsClient.GetAll(ctx, q, &pData); err != nil {
		return nil, err
	}
	return pData, nil
}
