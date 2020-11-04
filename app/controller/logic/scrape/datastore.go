package scrape

import (
	"context"

	"cloud.google.com/go/datastore"
	ds "github.com/niko-cb/japancovid19line/app/config/datastore"
	"github.com/niko-cb/japancovid19line/app/models"
)

// updateDatastore updates Datastore with the new Prefecture and Date data
func updateDatastore(ctx context.Context, data []*models.PrefectureData, date string) error {
	c, err := ds.NewClient()
	if err != nil {
		return err
	}

	if err := updatePrefectureData(ctx, c, data); err != nil {
		return err
	}

	if err := updateSourceDate(ctx, c, date); err != nil {
		return err
	}

	return nil
}

// updatePrefectureData updates the data for each Prefecture
func updatePrefectureData(ctx context.Context, c *datastore.Client, prefectureData []*models.PrefectureData) error {
	var keys []*datastore.Key
	for _, jpd := range prefectureData {
		key := datastore.NameKey(ds.DataKind(), jpd.Prefecture, nil)
		keys = append(keys, key)
	}

	if _, err := c.PutMulti(ctx, keys, prefectureData); err != nil {
		return err
	}

	return nil
}

// updateSourceDate updates the latest date
func updateSourceDate(ctx context.Context, c *datastore.Client, date string) error {
	sourceDate := new(models.SourceDate)
	sourceDate.Date = date

	dateKey := datastore.NameKey(ds.DateKind(), ds.DateName(), nil)
	if _, err := c.Put(ctx, dateKey, sourceDate); err != nil {
		return err
	}

	return nil
}
