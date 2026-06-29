package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

func makeMain(cars []Car, categories []Category, manufacturers []Manufacturer) []Main {

	var mains []Main
	for  _,  car := range cars {
		manName := "NO"
		manCountry := "NO"
		catName := "NO"
		for _, manufacturer := range manufacturers {
			if manufacturer.ID == car.ManufacturerID {
				manName = manufacturer.Name
				manCountry = manufacturer.Country
				break
			}
		}
		for _, category := range categories {
			if category.ID == car.CategoryID {
				catName = category.Name
				break
			}
		}
		main := Main{
			ID: car.ID,
			Name: car.Name,
			Year: car.Year,
			Image: car.Image,
			ManufacturerName: manName,
			ManufacturerCountry: manCountry,
			CategoryName: catName,

			Engine: car.Specifications.Engine,
			Horsepower: car.Specifications.Horsepower,
			Transmission: car.Specifications.Transmission,
			Drivetrain: car.Specifications.Drivetrain,
		}
		mains = append(mains, main)
	}
	return mains
}

func getCars() ([]Car, error) {
	data, err := fetchAPI("/models/")
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

func getCategories() ([]Category, error) {
	data, err := fetchAPI("/categories")
	if err != nil {
		return nil, err
	}

	var categories []Category

	err = json.Unmarshal(data, &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func getManufacturers() ([]Manufacturer, error) {
	data, err := fetchAPI("/manufacturers")
	if err != nil {
		return nil, err
	}

	var manufacturers []Manufacturer

	err = json.Unmarshal(data, &manufacturers)
	if err != nil {
		return nil, err
	}

	return manufacturers, nil
}


func getCar(id string) (Car, error) {
	data, err := fetchAPI("/models/" + id)

	var car Car

	if err != nil {
		return car, err
	}

	
	err = json.Unmarshal(data, &car)
	if err != nil {
		return car, err
	}

	return car, nil
}