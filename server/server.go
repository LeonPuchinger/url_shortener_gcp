package main

import (
	"context"
	"fmt"
	"net/http"
)

const Port = ":8080"

// Open firestore client and setup API routes.
func main() {
	ctx := context.Background()
	client, err := InitClient(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[len("/get/"):]
		url, err := GetUrl(ctx, client, key)
		if err != nil {
			w.WriteHeader(404)
			fmt.Fprint(w, err.Error())
			return
		}
		fmt.Fprint(w, url)
	})
	http.ListenAndServe(Port, nil)
}
