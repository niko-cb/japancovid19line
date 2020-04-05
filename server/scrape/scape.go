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
	covidDataURL = "https://www.mhlw.go.jp/stf/newpage_10688.html"
	// Element for the table that contains the japanese covid data table
	dataTableElementSelector = "#content > div.l-contentBody > div > div.l-contentMain > div:nth-child(4) > div > table:nth-child(106)"
	// element for each specific row in the japanese covid data table
	dataElementSelector = "tbody > tr > td"
	// element for the current data's date of publication
	dataSourceDateSelector = "#content > div.l-contentBody > div > div.l-contentMain > div:nth-child(4) > div"
)

func Scrape() []*model.PrefectureData {
	c := colly.NewCollector()
	var fullCovidArray []string
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	// gets table data
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

	var date string
	c.OnHTML(dataSourceDateSelector, func(e *colly.HTMLElement) {
		date = e.ChildText("u:nth-child(50)")
	})

	c.Visit(covidDataURL)
	return formatData(fullCovidArray, date)
}

func formatData(data []string, date string) []*model.PrefectureData {
	var dataSlice []*model.PrefectureData
	for i := 5; i < len(data); i += 5 {
		covidData := model.NewPrefectureData(data[i], data[i+1], data[i+2], data[i+3], data[i+4])
		dataSlice = append(dataSlice, covidData)
	}
	insertOrReinsertToDatastore(dataSlice, date)
	return dataSlice
}

func insertOrReinsertToDatastore(data []*model.PrefectureData, date string) {
	ctx := context.Background()
	dsClient, err := utils.NewDSClient()
	if err != nil {
		return
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
