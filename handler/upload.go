package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"appstore/model"
	"appstore/service"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse from body of request to get a json object.
	fmt.Println("Received one upload request")
	decoder := json.NewDecoder(r.Body)
	var app model.App
	if err := decoder.Decode(&app); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "Upload request received: %s\n", app.Description)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one search request")
	w.Header().Set("Content-Type", "application/json")
	title := r.URL.Query().Get("title")
	description := r.URL.Query().Get("description")

	var apps []model.App
	var err error
	apps, err = service.SearchApps(title, description)
	if err != nil {
		http.Error(w, "Failed to read post from backend", http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(apps)
	if err != nil {
		http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
