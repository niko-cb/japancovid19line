package model

import (
	"context"
	"strconv"

	"cloud.google.com/go/datastore"
)

type PrefectureData struct {
	Prefecture string `json:"prefecture"`
	Cases      string `json:"cases"`
	Recovered  string `json:"recovered"`
	Deaths     string `json:"deaths"`
}

type PrefectureDataSlice = []*PrefectureData

func NewPrefectureData(pref string, cases, rec, deaths int) *PrefectureData {
	return &PrefectureData{
		Prefecture: pref,
		Cases:      strconv.Itoa(cases),
		Recovered:  strconv.Itoa(rec),
		Deaths:     strconv.Itoa(deaths),
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
