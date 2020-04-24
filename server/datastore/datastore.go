package datastore

import (
	"os"

	"github.com/niko-cb/covid19datascraper/server/context"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

func NewClient() (*datastore.Client, error) {
	ctx := context.Get()
	dsJSON := os.Getenv("DATASTORE_CRED_JSON")
	projectID := os.Getenv("PROJECT_ID")
	if c, err := datastore.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(dsJSON))); err != nil {
		return nil, err
	} else {
		return c, nil
	}
}

func DataKind() string {
	return "CovidData"
}

func DateKind() string {
	return "SourceDate"
}
