package scrape

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/niko-cb/covid19datascraper/server/model"
	"github.com/niko-cb/covid19datascraper/server/utils"
	"log"
)

const (
	// Article URL for current day
	covidDataURL = "https://www.mhlw.go.jp/stf/newpage_10651.html"
	// Element for the table that contains the japanese covid data table
	dataTableElementSelector = "#content > div.l-contentBody > div > div.l-contentMain > div:nth-child(4) > div > table:nth-child(107)"
	// element for each specific row in the japanese covid data table
	dataElementSelector = "tbody > tr > td.aly_tx_center "
)

func Scrape() []*model.PrefectureData {
	c := colly.NewCollector()
	var fullCovidArray []string

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML(dataTableElementSelector, func(e *colly.HTMLElement) {
		e.ForEach(dataElementSelector, func(_ int, e *colly.HTMLElement) {
			var data string
			data = e.ChildText("font")
			if data == "" {
				return
			}
			fullCovidArray = append(fullCovidArray, data)
		})
	})

	c.Visit(covidDataURL)

	return dataToModel(fullCovidArray)
}

func dataToModel(data []string) []*model.PrefectureData {
	var dataSlice []*model.PrefectureData
	for i := 5; i < len(data); i += 5 {
		covidData := model.NewPrefectureData(data[i], data[i+1], data[i+2], data[i+3], data[i+4])
		dataSlice = append(dataSlice, covidData)
	}
	insertOrReinsertToDatastore(dataSlice)
	return dataSlice
}

func insertOrReinsertToDatastore(data []*model.PrefectureData) {
	ctx := context.Background()
	dsClient, err := utils.NewDSClient()
	if err != nil {
		return
	}

	kind := utils.DatastoreKind()

	for _, jpd := range data {
		name := jpd.Prefecture
		key := datastore.NameKey(kind, name, nil)

		// Saves the new entity.
		if _, err := dsClient.Put(ctx, key, jpd); err != nil {
			log.Fatalf("Failed to save data: %v", err)
		}

		fmt.Printf("saved %s data\n", name)
	}
}
