package model

import (
	"cloud.google.com/go/datastore"
	"context"
)

type SourceDate struct {
	Date string
}

func GetSourceDate(ctx context.Context, dsClient *datastore.Client, kind string) (*SourceDate, error) {
	sd := new(SourceDate)
	key := datastore.NameKey(kind, "Latest", nil)

	if err := dsClient.Get(ctx, key, sd); err != nil {
		return nil, err
	}
	return sd, nil
}
