package scrape

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/PuerkitoBio/goquery"
	"github.com/niko-cb/covid19datascraper/server/model"
	"github.com/niko-cb/covid19datascraper/server/utils"
)

const (
	// Check Latest JSON File
	covidDataJSONLatestURL = "https://raw.githubusercontent.com/reustle/covid19japan-data/master/docs/summary/latest.json"
	// URL for latest JSON File
	covidDataJSON = "https://raw.githubusercontent.com/reustle/covid19japan-data/master/docs/summary/"
)

func Scrape() []*model.PrefectureData {
	latest := getLatestJson()
	data, err := readJSONFromUrl(covidDataJSON + latest)
	if err != nil {
		log.Println(err.Error())
	}
	date := latest[:10]
	var pData []*model.PrefectureData
	for _, city := range data.Prefectures {
		prefectureData := model.NewPrefectureData(city.NameJA, city.Confirmed, city.Deaths, city.Recovered)
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
	dsClient, err := utils.NewDSClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	kind := utils.DatastoreKind()

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

	dateKind := utils.DatastoreDateKind()
	name := "Latest"
	dateKey := datastore.NameKey(dateKind, name, nil)
	if _, err := dsClient.Put(ctx, dateKey, sourceDate); err != nil {
		log.Fatalf("failed to save date into datastore: %v", err)
	}
}
