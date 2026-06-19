package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

const carsAPI = "http://localhost:3000/api"

type Specification struct {
	Engine       string `json:"engine"`
	Horsepower   int    `json:"horsepower"`
	Transmission string `json:"transmission"`
	Drivetrain   string `json:"drivetrain"`
}

type Car struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	ManufacturerID int           `json:"manufacturerId"`
	CategoryID     int           `json:"categoryId"`
	Year           int           `json:"year"`
	Specifications Specification `json:"specifications"`
	Image          string        `json:"image"`
}

func fetchAPI(endpoint string) ([]byte, error) {
	response, err := http.Get(carsAPI + endpoint)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getCars() ([]Car, error) {
	data, err := fetchAPI("/models")
	if err != nil {
		return nil, err
	}

	var cars []Car

	err = json.Unmarshal(data, &cars)
	if err != nil {
		return nil, err
	}

	return cars, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	cars, err := getCars()
	if err != nil {
		http.Error(w, "Server issue: could not load cars right now.", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Server issue: could not load page template.", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, cars)
}


func main() {
	http.HandleFunc("/", homeHandler)
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}