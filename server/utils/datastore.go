package utils

import (
	"cloud.google.com/go/datastore"
	"context"
	"google.golang.org/api/option"
)

func NewDSClient() (*datastore.Client, error) {
	ctx := context.Background()
	if c, err := datastore.NewClient(ctx, "japancovid19", option.WithCredentialsFile("datastorecredentials.json")); err != nil {
		return nil, err
	} else {
		return c, nil
	}
}
