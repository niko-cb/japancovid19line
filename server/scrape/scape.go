package scrape

import (
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/niko-cb/covid19datascraper/server/model"
	"github.com/niko-cb/covid19datascraper/server/utils"
	"log"
	"net/http"
)

const (
	// Check Latest JSON File
	covidDataJSONLatestURL = "https://raw.githubusercontent.com/reustle/covid19japan-data/master/docs/summary/latest.json"
	// URL for latest JSON File
	covidDataJSON = "https://raw.githubusercontent.com/reustle/covid19japan-data/master/docs/summary/"
)

func Scrape() []*model.PrefectureData {
	resp, err := http.Get(covidDataJSONLatestURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	q, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	latestJsonFile := q.Text()
	fileUrl := covidDataJSON + latestJsonFile

	data, err := readJSONFromUrl(fileUrl)
	if err != nil {
		log.Println(err.Error())
	}

	date := latestJsonFile[:10]
	log.Println(date)

	var pData []*model.PrefectureData

	for _, city := range data.Prefectures {
		prefectureData := new(model.PrefectureData)
		prefectureData.Prefecture = city.Name
		prefectureData.Cases = string(city.Confirmed)
		prefectureData.Deaths = string(city.Deaths)
		prefectureData.Recovered = string(city.Recovered)
	}

	insertOrReinsertToDatastore(pData, date)
	return pData

}

func readJSONFromUrl(url string) (*model.AllData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var pData *model.AllData
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()
	if err := json.Unmarshal(respByte, &pData); err != nil {
		return nil, err
	}

	return pData, nil
}

func insertOrReinsertToDatastore(data []*model.PrefectureData, date string) {
	ctx := context.Background()
	dsClient, err := utils.NewDSClient()
	if err != nil {
		log.Fatalf(err.Error())
	}

	kind := utils.DatastoreKind()

	log.Println(date)

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
