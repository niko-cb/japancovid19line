package scrape

import (
	"encoding/json"

	"github.com/niko-cb/covid19datascraper/app/models"
	"github.com/niko-cb/covid19datascraper/app/models/helpers/prefectures"
)

// formatData takes the response and turns it into a Prefecture Data struct slice
func formatData(bytes []byte) ([]*models.PrefectureData, error) {
	var res *models.CovidDataRes
	if err := json.Unmarshal(bytes, &res); err != nil {
		return nil, err
	}

	var pData []*models.PrefectureData
	for _, p := range res.Prefectures {
		if p.NameJA == "" {
			p.NameJA = (&prefectures.PrefectureMap{}).Japanese(p.Name)
		}
		cities, err := json.Marshal(p.ConfirmedByCity)
		if err != nil {
			return nil, err
		}

		prefectureData := models.NewPrefectureData(p.NameJA, p.Confirmed,
			p.Deaths, p.Recovered, p.NewlyConfirmed,
			p.YesterdayConfirmed, string(cities))

		pData = append(pData, prefectureData)
	}
	return pData, nil
}
