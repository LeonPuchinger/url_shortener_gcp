package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const Port = ":8080"
const ErrorPage = "https://storage.googleapis.com/static-client/error.html"

//allow CORS from any origin
func allowCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func setContentTypeJSON(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
}

//encode a map as json and serves it with a status-code
type jsonHandler func(http.ResponseWriter, *http.Request) (int, map[string]string)

func (handler jsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	allowCORS(&w)
	setContentTypeJSON(&w)
	status, body := handler(w, r)
	response, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
		return
	}
	w.WriteHeader(status)
	w.Write(response)
}

// Open firestore client and setup API routes.
func main() {
	ctx := context.Background()
	client, err := InitClient(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	urlReachable := func(url string) bool {
		_, err := http.Get(url)
		return err == nil
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		allowCORS(&w)
		if r.Method == http.MethodGet {
			key := r.URL.Path[len("/"):]
			url, err := GetUrl(ctx, client, key)
			if err != nil {
				http.Redirect(w, r, ErrorPage, http.StatusFound)
				return
			}
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		w.WriteHeader(405)
	})
	http.Handle("/get/", jsonHandler(func(w http.ResponseWriter, r *http.Request) (int, map[string]string) {
		if r.Method == http.MethodGet {
			key := r.URL.Path[len("/get/"):]
			url, err := GetUrl(ctx, client, key)
			if err != nil {
				return 404, map[string]string{"error": "could not get url for this key"}
			}
			return 200, map[string]string{
				"key": key,
				"url": url,
			}
		}
		return 405, map[string]string{"error": "method not allowed"}
	}))
	http.Handle("/add", jsonHandler(func(w http.ResponseWriter, r *http.Request) (int, map[string]string) {
		if r.Method == http.MethodPost {
			url := r.FormValue("url")
			if url == "" {
				return 400, map[string]string{"error": "url form parameter missing or empty"}
			}
			if !urlReachable(url) {
				return 400, map[string]string{"error": "url is not reachable"}
			}
			key, err := AddUrl(ctx, client, url)
			if err != nil {
				return 500, map[string]string{"error": "unable to add url"}
			}
			return 200, map[string]string{
				"key": key,
				"url": url,
			}
		}
		return 405, map[string]string{"error": "method not allowed"}
	}))
	err = http.ListenAndServe(Port, nil)
	if err != nil {
		fmt.Printf("could not open http server: %s\n", err.Error())
	}
}
