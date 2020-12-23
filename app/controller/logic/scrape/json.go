package scrape

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// readJSONFromLatestURL gets the latest json data from the latest url
func readJSONFromLatestURL() (string, error) {
	res, err := http.Get(covidDataJSONLatestURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	d, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}
	return d.Text(), nil
}
