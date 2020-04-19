package model

import (
	"cloud.google.com/go/datastore"
	"context"
)

type PrefectureData struct {
	Prefecture string `json:"prefecture"`
	Cases      int    `json:"cases"`
	Recovered  int    `json:"recovered"`
	Deaths     int    `json:"deaths"`
}

type PrefectureDataSlice = []*PrefectureData

func NewPrefectureData(pref string, cases, rec, deaths int) *PrefectureData {
	return &PrefectureData{
		Prefecture: pref,
		Cases:      cases,
		Recovered:  rec,
		Deaths:     deaths,
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
