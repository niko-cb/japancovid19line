package utils

import (
	"cloud.google.com/go/datastore"
	"context"
	"google.golang.org/api/option"
	"os"
)

func NewDSClient() (*datastore.Client, error) {
	ctx := context.Background()
	dsJSON := os.Getenv("DATASTORE_CRED_JSON")
	if c, err := datastore.NewClient(ctx, "japancovid19", option.WithCredentialsJSON([]byte(dsJSON))); err != nil {
		return nil, err
	} else {
		return c, nil
	}
}

func DatastoreKind() string {
	return "CovidData"
}

func DatastoreDateKind() string {
	return "SourceDate"
}
