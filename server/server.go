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

	allowCORS := func(w *http.ResponseWriter) {
		(*w).Header().Set("Access-Control-Allow-Origin", "*")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		allowCORS(&w)
		w.WriteHeader(404)
	})
	http.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		allowCORS(&w)
		if r.Method == http.MethodGet {
			key := r.URL.Path[len("/get/"):]
			url, err := GetUrl(ctx, client, key)
			if err != nil {
				w.WriteHeader(404)
				fmt.Fprint(w, err.Error())
				return
			}
			fmt.Fprint(w, url)
		} else {
			w.WriteHeader(405)
		}
	})
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		allowCORS(&w)
		if r.Method == http.MethodPost {
			url := r.FormValue("url")
			if url == "" {
				w.WriteHeader(400)
				fmt.Fprint(w, "url form parameter missing or empty")
				return
			}
			key, err := AddUrl(ctx, client, url)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprint(w, err.Error())
				return
			}
			fmt.Fprint(w, key)
		} else {
			w.WriteHeader(405)
		}
	})
	err = http.ListenAndServe(Port, nil)
	if err != nil {
		fmt.Printf("could not open http server: %s\n", err.Error())
	}
}
