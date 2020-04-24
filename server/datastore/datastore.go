package datastore

import (
	"github.com/niko-cb/covid19datascraper/server/env"

	"github.com/niko-cb/covid19datascraper/server/context"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

func NewClient() (*datastore.Client, error) {
	ctx := context.Get()
	dsJSON := env.AuthDatastore()
	projectID := env.ProjectID()
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
