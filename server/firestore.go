package main

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
)

const ProjectId = "url-shortener-308812"

// Opens a new firestore client.
// This usually fails if there is invalid or no authentication provided
func InitClient(ctx context.Context) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, ProjectId)
	if err != nil {
		return nil, errors.New("could not open firestore client")
	}
	return client, nil
}

// Queries url from firestore database by key.
// Returns the url string or an error if the url could not be found with the provided key
func GetUrl(ctx context.Context, client *firestore.Client, key string) (string, error) {
	query := client.Collection("urls").Where("key", "==", key)
	doc, err := query.Documents(ctx).Next()
	if err != nil {
		return "", fmt.Errorf(`no url found for key "%s"`, key)
	}
	type Url struct {
		Url string `firestore:"url"`
	}
	var url Url
	err = doc.DataTo(&url)
	if err != nil || url.Url == "" {
		return "", errors.New("could not get url from queried document")
	}
	return url.Url, nil
}
