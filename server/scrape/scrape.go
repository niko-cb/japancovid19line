package scrape

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	datastore2 "github.com/niko-cb/covid19datascraper/server/datastore"

	"github.com/niko-cb/covid19datascraper/server/prefectures"

	"cloud.google.com/go/datastore"
	"github.com/PuerkitoBio/goquery"
	"github.com/niko-cb/covid19datascraper/server/model"
)

const (
	// Check Latest JSON File
	covidDataJSONLatestURL = "https://raw.githubusercontent.com/reustle/covid19japan-data/master/docs/summary/latest.json"
	// URL for latest JSON File
	covidDataJSON = "https://raw.githubusercontent.com/reustle/covid19japan-data/master/docs/summary/"
)

func Do() []*model.PrefectureData {
	latest := getLatestJson()
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
		prefectureData := model.NewPrefectureData(pref, prefecture.Confirmed, prefecture.Deaths, prefecture.Recovered, string(cities))
		pData = append(pData, prefectureData)
	}
	updateDatastore(pData, date)
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

func getLatestJson() string {
	resp, err := http.Get(covidDataJSONLatestURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	q, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	return q.Text()
}

func updateDatastore(data []*model.PrefectureData, date string) {
	ctx := context.Background()
	dsClient, err := datastore2.NewClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	kind := datastore2.DataKind()

	var keys []*datastore.Key
	for _, jpd := range data {
		name := jpd.Prefecture
		key := datastore.NameKey(kind, name, nil)
		keys = append(keys, key)
	}

	if _, err := dsClient.PutMulti(ctx, keys, data); err != nil {
		log.Fatalf("Failed to save data: %v", err)
	}

	sourceDate := new(model.SourceDate)
	sourceDate.Date = date

	dateKind := datastore2.DateKind()
	name := "Latest"
	dateKey := datastore.NameKey(dateKind, name, nil)
	if _, err := dsClient.Put(ctx, dateKey, sourceDate); err != nil {
		log.Fatalf("failed to save date into datastore: %v", err)
	}
}
