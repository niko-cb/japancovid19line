package datastore

import (
	"github.com/niko-cb/covid19datascraper/server/context"
	"github.com/niko-cb/covid19datascraper/server/env"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

func NewClient() (*datastore.Client, error) {
	ctx := context.Get()
	config := env.Get()
	dsJSON := config.DatastoreAuth
	projectID := config.ProjectID
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

func DateName() string {
	return "Latest"
}
