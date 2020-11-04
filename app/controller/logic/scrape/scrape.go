package scrape

import (
	"context"
)

func Do(ctx context.Context) error {
	latest, err := getLatestJson()
	if err != nil {
		return err
	}

	dataBytes, err := readJSONFromURL(covidDataJSON + latest)
	if err != nil {
		return err
	}

	prefData, err := formatData(dataBytes)
	if err != nil {
		return err
	}

	date := latest[:10]
	if err := updateDatastore(ctx, prefData, date); err != nil {
		return err
	}

	return nil
}
