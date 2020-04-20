package utils

import (
	"os"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

func NewDSClient() (*datastore.Client, error) {
	ctx := GetContext()
	dsJSON := os.Getenv("DATASTORE_CRED_JSON")
	projectID := os.Getenv("PROJECT_ID")
	if c, err := datastore.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(dsJSON))); err != nil {
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
