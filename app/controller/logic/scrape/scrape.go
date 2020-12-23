package scrape

import (
	"context"
	"strings"
)

func Do(ctx context.Context) error {
	latest, err := getLatestJson()
	if err != nil {
		return err
	}

	prefData, date, err := formatData([]byte(latest))
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
