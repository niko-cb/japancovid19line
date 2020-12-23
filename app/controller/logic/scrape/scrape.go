package scrape

import (
	"context"
	"strings"
)

func Do(ctx context.Context) error {
	latestDataString, err := readJSONFromLatestURL()
	if err != nil {
		return err
	}

	prefData, date, err := formatData([]byte(latestDataString))
	if err != nil {
		return err
	}

	if date != "" {
		date = strings.Split(date, "T")[0]
	}

	if err := updateDatastore(ctx, prefData, date); err != nil {
		return err
	}

	return nil
}
