package scrape

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/niko-cb/covid19datascraper/server/model"
	"github.com/niko-cb/covid19datascraper/server/utils"
	"log"
	"time"
)

const (
	// Article URL for current day
	covidDataURL = "https://docs.google.com/spreadsheets/d/e/2PACX-1vRri4r42DHwMHePjJfYN-qEWhGvKeOQullBtEzfle15i-xAsm9ZgV8oMxQNhPRO1CId39BPnn1IO5YO/pubhtml#"
	// Element for the table that contains the japanese covid data table
	dataTableElementSelector = "#1399411442 > div > table > tbody"
	// element for each specific row in the japanese covid data table
	dataElementSelector = "tr"
	// element for the current data's date of publication
	dataSourceDateSelector = "#434907293 > div > table > tbody > tr:nth-child(20)"
)

func Scrape() []*model.PrefectureData {
	c := colly.NewCollector()
	var covidData []string

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.SetRequestTimeout(5 * time.Minute)

	// gets table data
	c.OnHTML(dataTableElementSelector, func(e *colly.HTMLElement) {
		e.ForEach(dataElementSelector, func(_ int, e *colly.HTMLElement) {
			var pref string
			var cases string
			var rec string
			var deaths string
			pref = e.ChildText("td:nth-child(3)")
			cases = e.ChildText("td:nth-child(4)")
			if cases == "" {
				cases = "0"
			}
			rec = e.ChildText("td:nth-child(5)")
			if rec == "" {
				rec = "0"
			}
			deaths = e.ChildText("td:nth-child(6)")
			if deaths == "" {
				deaths = "0"
			}
			covidData = append(covidData, pref)
			covidData = append(covidData, cases)
			covidData = append(covidData, rec)
			covidData = append(covidData, deaths)
		})
	})

	var date string
	c.OnHTML(dataSourceDateSelector, func(e *colly.HTMLElement) {
		date = e.ChildText("td:nth-child(6)")
	})

	err := c.Visit(covidDataURL)
	if err != nil {
		log.Println(err.Error())
	}

	return formatData(covidData, date)
}

func formatData(data []string, date string) []*model.PrefectureData {
	var dataSlice []*model.PrefectureData
	// Skip the headers for the table (start i at 12)
	for i := 12; i < len(data); i += 4 {
		covidData := model.NewPrefectureData(data[i], data[i+1], data[i+2], data[i+3])
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
