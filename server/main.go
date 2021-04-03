package main

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
)

func main() {
	projectId := "url-shortener-308812"
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		fmt.Println("ERROR: could not open client")
		fmt.Println(err)
		return
	}

	url, err := getUrl(ctx, client, "0") //example query
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	fmt.Println(url)
}

// Queries url from firestore database by key.
// Returns the url string or an error if the url could not be found with the provided key
func getUrl(ctx context.Context, client *firestore.Client, key string) (string, error) {
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
