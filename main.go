package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

const carsAPI = "http://localhost:3000/api"

func homeHandler(w http.ResponseWriter, r *http.Request) {
	result := make(chan dataResult) 
	go func() {
		cars, err := getCars()
		result <-  dataResult{
			Name: "cars",
			Cars:  cars,
			Error: err,
		}
	}()
	go func() {
		categories, err := getCategories()
		result <-  dataResult{
			Name: "categories",
			Categories:  categories,
			Error: err,
		}
	}()
	go func() {
		manufacturers, err := getManufacturers()
		result <-  dataResult{
			Name: "manufacturers",
			Manufacturers:  manufacturers,
			Error: err,
		}
	}()
	var  cars []Car
	var categories []Category
	var manufacturers []Manufacturer
	for i := 0; i < 3; i++ {
		finalData :=  <-result
		if finalData.Error != nil {
			http.Error(w, "Server issue: could not load categories right now.", http.StatusInternalServerError)
			return
		}
		if finalData.Name == "cars" {
			cars = finalData.Cars
		}
		if finalData.Name == "categories" {
			categories = finalData.Categories
		}
		if finalData.Name == "manufacturers" {
			manufacturers = finalData.Manufacturers
		}
	}
	mains := makeMain(cars, categories, manufacturers)
	filters := FilterOptions{
		Search:        r.URL.Query().Get("search"),
		Manufacturer:  r.URL.Query().Get("manufacturer"),
		Country:       r.URL.Query().Get("country"),
		Category:      r.URL.Query().Get("category"),
		MinYear:       r.URL.Query().Get("minYear"),
		Engine:        r.URL.Query().Get("engine"),
		MinHorsepower: r.URL.Query().Get("minHorsepower"),
		Transmission:  r.URL.Query().Get("transmission"),
		Drivetrain:    r.URL.Query().Get("drivetrain"),
	}
	//mains = filterCarCardsByName(mains, search)
	optionsData := getUniqueOptions(mains)
	filteredCards := filterMainCards(mains, filters)
	pageData := PageData{
		Filters: filters,

		ManufacturerOptions: optionsData.ManufacturerOptions,
		CountryOptions:      optionsData.CountryOptions,
		CategoryOptions:     optionsData.CategoryOptions,
		EngineOptions:       optionsData.EngineOptions,
		TransmissionOptions: optionsData.TransmissionOptions,
		DrivetrainOptions:   optionsData.DrivetrainOptions,

		Mains: filteredCards,
	}
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Server issue: page template could not be loaded.", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, "Server issue: page could not be rendered.", http.StatusInternalServerError)
		return
	}
}

func specsHandler(w http.ResponseWriter, r *http.Request){
	
	url := r.URL.Path
	parts := strings.Split(strings.TrimPrefix(url, "/"), "/")
	id :=  parts[1]
	car, err := getCar(id)
	if err != nil {
		http.Error(w, "Server issue: could not load cars right now.", http.StatusInternalServerError)
			return
	}

		tmpl, err := template.ParseFiles("specs.html")
		tmpl.Execute(w,car)
}


func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", homeHandler)
	mux.HandleFunc("GET /specifications/", specsHandler)
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}