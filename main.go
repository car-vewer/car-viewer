package main

import (
	"html/template"
	"log"
	"net/http"
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
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Server issue: could not load page template.", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, mains)
}


func main() {
	http.HandleFunc("/", homeHandler)
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}