package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"strconv"
)

const carsAPI = "http://localhost:3000/api"
var comparisonIDs []string


func containsID(ids []string, id string) bool {
	for _, existingID := range ids {
		if existingID == id {
			return true
		}
	}

	return false
}

func addComparisonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	id := r.FormValue("id")

	if id != "" && !containsID(comparisonIDs, id) && len(comparisonIDs) < 2 {
		comparisonIDs = append(comparisonIDs, id)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func comparisonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	cars, categories, manufacturers, err := loadData()
	if err != nil {
		http.Error(w, "Server issue: could not load comparison data right now.", http.StatusInternalServerError)
		return
	}

	mains := makeMain(cars, categories, manufacturers)

	var comparedCars []Main

	for _, car := range mains {
		carID := strconv.Itoa(car.ID)

		if containsID(comparisonIDs, carID) {
			comparedCars = append(comparedCars, car)
		}
	}

	pageData := ComparisonPageData{
		Cars: comparedCars,
	}

	tmpl, err := template.ParseFiles("comparison.html")
	if err != nil {
		http.Error(w, "Server issue: comparison page could not be loaded.", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, "Server issue: comparison page could not be rendered.", http.StatusInternalServerError)
		return
	}
}

func removeComparisonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	id := r.FormValue("id")

	if id != "" {
		comparisonIDs = removeID(comparisonIDs, id)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func removeID(ids []string, id string) []string {
	var updated []string

	for _, existingID := range ids {
		if existingID != id {
			updated = append(updated, existingID)
		}
	}

	return updated
}

func markCompared(mains []Main) []Main {
	for i := range mains {
		carID := strconv.Itoa(mains[i].ID)

		if containsID(comparisonIDs, carID) {
			mains[i].IsCompared = true
		}
	}

	return mains
}

func loadData() (cars []Car, categories []Category, manufacturers []Manufacturer, err error) {
	result := make(chan dataResult)
	go func() {
		cars, err := getCars()
		result <- dataResult{
			Name:  "cars",
			Cars:  cars,
			Error: err,
		}
	}()
	go func() {
		categories, err := getCategories()
		result <- dataResult{
			Name:       "categories",
			Categories: categories,
			Error:      err,
		}
	}()
	go func() {
		manufacturers, err := getManufacturers()
		result <- dataResult{
			Name:          "manufacturers",
			Manufacturers: manufacturers,
			Error:         err,
		}
	}()
	for i := 0; i < 3; i++ {
		finalData := <-result
		if finalData.Error != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	http.Error(w, "Server issue: could not load categories right now.", http.StatusInternalServerError)
			err = finalData.Error
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
	return cars, categories, manufacturers, err
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	cars, categories, manufacturers, err := loadData()
	if err != nil {
		http.Error(w, "Server issue: could not load comparison data right now.", http.StatusInternalServerError)
		return
	}
	mains := makeMain(cars, categories, manufacturers)
	mains = markCompared(mains)
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
		CompareCount: len(comparisonIDs),
	}
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Server issue: page template could not be loaded.", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Server issue: page could not be rendered.", http.StatusInternalServerError)
		return
	}
	//w.WriteHeader(http.StatusOK)
}

func specsHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path
	parts := strings.Split(strings.TrimPrefix(url, "/"), "/")
	id := parts[1]
	car, err := getCar(id)
	if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Server issue: could not load car details for this index.", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("specs.html")
	err = tmpl.Execute(w, car)
	if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Server issue: page could not be rendered.", http.StatusInternalServerError)
		return
	}
	//w.WriteHeader(http.StatusOK)

}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", homeHandler)
	mux.HandleFunc("GET /specifications/", specsHandler)
	mux.HandleFunc("GET /comparison", comparisonHandler)
	mux.HandleFunc("POST /comparison/add", addComparisonHandler)
	mux.HandleFunc("POST /comparison/remove", removeComparisonHandler)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
