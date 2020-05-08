package scrape

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/PuerkitoBio/goquery"
	ds "github.com/niko-cb/covid19datascraper/server/datastore"
	"github.com/niko-cb/covid19datascraper/server/model"
	"github.com/niko-cb/covid19datascraper/server/prefectures"
)

func Do(ctx context.Context) []*model.PrefectureData {
	latest, err := getLatestJson()
	if err != nil {
		log.Fatal(err)
	}
	data, err := readJSONFromUrl(covidDataJSON + latest)
	if err != nil {
		log.Println(err.Error())
	}
	date := latest[:10]
	var pData []*model.PrefectureData
	for _, prefecture := range data.Prefectures {
		pref := prefecture.NameJA
		if prefecture.NameJA == "" {
			pref = (&prefectures.PrefectureMap{}).Japanese(prefecture.Name)
		}
		cities, err := json.Marshal(prefecture.ConfirmedByCity)
		if err != nil {
			log.Println(err.Error())
		}

		prefectureData := model.NewPrefectureData(pref, prefecture.Confirmed,
			prefecture.Deaths, prefecture.Recovered, prefecture.NewlyConfirmed,
			prefecture.YesterdayConfirmed, string(cities))

		pData = append(pData, prefectureData)
	}
	updateDatastore(ctx, pData, date)
	return pData
}

func readJSONFromUrl(url string) (*model.CovidDataRes, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var pData *model.CovidDataRes
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()
	if err := json.Unmarshal(respByte, &pData); err != nil {
		return nil, err
	}

	return pData, nil
}

func getLatestJson() (string, error) {
	resp, err := http.Get(covidDataJSONLatestURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	q, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	return q.Text(), nil
}

func updateDatastore(ctx context.Context, data []*model.PrefectureData, date string) {
	dsClient, err := ds.NewClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	var keys []*datastore.Key
	for _, jpd := range data {
		key := datastore.NameKey(ds.DataKind(), jpd.Prefecture, nil)
		keys = append(keys, key)
	}

	if _, err := dsClient.PutMulti(ctx, keys, data); err != nil {
		log.Fatalf("failed to save data into datastore: %v", err)
	}

	sourceDate := new(model.SourceDate)
	sourceDate.Date = date

	dateKey := datastore.NameKey(ds.DataKind(), ds.DateName(), nil)
	if _, err := dsClient.Put(ctx, dateKey, sourceDate); err != nil {
		log.Fatalf("failed to save date into datastore: %v", err)
	}
}
