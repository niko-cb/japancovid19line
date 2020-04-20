package model

import (
	"context"

	"cloud.google.com/go/datastore"
)

type SourceDate struct {
	Date string
}

func GetDateFromDatastore(ctx context.Context, dsClient *datastore.Client, kind string) (*SourceDate, error) {
	sd := new(SourceDate)
	key := datastore.NameKey(kind, "Latest", nil)

	if err := dsClient.Get(ctx, key, sd); err != nil {
		return nil, err
	}
	return sd, nil
}
