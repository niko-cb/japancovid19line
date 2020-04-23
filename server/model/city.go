package model

import (
	"context"
	"strconv"

	"cloud.google.com/go/datastore"
)

type CityData struct {
	City  string `json:"city"`
	Cases string `json:"cases"`
}

type CityDataSlice = []*CityData

func NewCityData(city string, cases int) *CityData {
	return &CityData{
		City:  city,
		Cases: strconv.Itoa(cases),
	}
}

func GetCityData(ctx context.Context, dsClient *datastore.Client, kind string) ([]*PrefectureData, error) {
	var pData []*PrefectureData
	q := datastore.NewQuery(kind)
	if _, err := dsClient.GetAll(ctx, q, &pData); err != nil {
		return nil, err
	}
	return pData, nil
}
