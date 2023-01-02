package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"appstore/model"
	"appstore/service"

	"github.com/pborman/uuid"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse from body of request to get a json object.
	fmt.Println("Received one upload request")

	app := model.App{
		Id:          uuid.New(),
		User:        r.FormValue("user"),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	fmt.Printf("%v,%T", price, price)
	if err != nil {
		fmt.Println(err)
	}
	app.Price = int(price * 100.0)

	file, _, err := r.FormFile("media_file")
	if err != nil {
		http.Error(w, "Media file is not available", http.StatusBadRequest)
		fmt.Printf("Media file is not available %v\n", err)
		return
	}

	service.SaveApp(&app, file)
	if err != nil {
		http.Error(w, "Failed to save app to backend", http.StatusInternalServerError)
		fmt.Printf("Failed to save app to backend %v\n", err)
		return
	}

	fmt.Fprintf(w, "App is saved successfully: %s\n", app.Description)
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
		http.Error(w, "Failed to read App from backend", http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(apps)
	if err != nil {
		http.Error(w, "Failed to parse Apps into JSON format", http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
