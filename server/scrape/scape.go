package scrape

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/niko-cb/covid19datascraper/server/model"
	"github.com/niko-cb/covid19datascraper/server/utils"
	"log"
	"net/http"
)

const (
	// Article URL for current day
	covidDataURL = "https://docs.google.com/spreadsheets/d/e/2PACX-1vRri4r42DHwMHePjJfYN-qEWhGvKeOQullBtEzfle15i-xAsm9ZgV8oMxQNhPRO1CId39BPnn1IO5YO/pubhtml"
	// Element for the table that contains the japanese covid data table
	dataTableElementSelector = "#1399411442 > div > table > tbody"
	// element for each specific row in the japanese covid data table
	dataElementSelector = "tr"
	// element for the current data's date of publication
	dataSourceDateSelector = "#434907293 > div > table > tbody > tr:nth-child(20)"
)

func Scrape() []*model.PrefectureData {
	var covidData []string

	resp, err := http.Get(covidDataURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	q, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	q.Find(dataTableElementSelector).Each(func(index int, tableHtml *goquery.Selection) {
		tableHtml.Find(dataElementSelector).Each(func(index int, rowHtml *goquery.Selection) {
			rowHtml.Find("td:nth-child(3)").Each(func(index int, prefHtml *goquery.Selection) {
				var pref string
				pref = prefHtml.Text()
				covidData = append(covidData, pref)
			})
			rowHtml.Find("td:nth-child(4)").Each(func(index int, caseHtml *goquery.Selection) {
				var cases string
				cases = caseHtml.Text()
				if cases == "" {
					cases = "0"
				}
				covidData = append(covidData, cases)
			})
			rowHtml.Find("td:nth-child(5)").Each(func(index int, RecHtml *goquery.Selection) {
				var rec string
				rec = RecHtml.Text()
				if rec == "" {
					rec = "0"
				}
				covidData = append(covidData, rec)
			})
			rowHtml.Find("td:nth-child(6)").Each(func(index int, deathsHtml *goquery.Selection) {
				var deaths string
				deaths = deathsHtml.Text()
				if deaths == "" {
					deaths = "0"
				}
				covidData = append(covidData, deaths)
			})
		})
	})
	var date string
	q.Find(dataSourceDateSelector).Each(func(index int, dateHtml *goquery.Selection) {
		date = dateHtml.Find("td:nth-child(6)").Text()
	})
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
