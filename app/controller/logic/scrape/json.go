package scrape

import (
	"bytes"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// getLatestJson gets the latest json data from the latest url
func getLatestJson() (string, error) {
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

// readJSONFromURL gets the json data from the specified URL and returns it
func readJSONFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)

	return buf.Bytes(), nil
}
