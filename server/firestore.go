package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

// Increments a key by one. A key consists of multiple characters out of <0-9A-Za-z> 
// with 0 being the lowest and z the highest individual char.
func incrementKey(key string) string {
	len := len(key)
	if len == 0 {
		return "0"
	}
	builder := strings.Builder{}
	carry := true
	for indx := 0; indx < len; indx++ {
		char := key[indx]
		if carry {
			if char == '9' {
				builder.WriteRune('A')
				carry = false
			} else if char == 'Z' {
				builder.WriteRune('a')
				carry = false
			} else if char == 'z' {
				builder.WriteRune('0')
			} else {
				builder.WriteRune(rune(char + 1))
				carry = false
			}
			continue
		}
		builder.WriteByte(char)
	}
	if carry {
		builder.WriteRune('1')
	}
	return builder.String()
}
