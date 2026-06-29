package main

import (
	"strings"
	"strconv"
)

func getUniqueOptions(mainCards []Main) PageData {
	var pageData PageData

	manufacturerMap := make(map[string]bool)
	countryMap := make(map[string]bool)
	categoryMap := make(map[string]bool)
	engineMap := make(map[string]bool)
	transmissionMap := make(map[string]bool)
	drivetrainMap := make(map[string]bool)

	for _, car := range mainCards {
		if car.ManufacturerName != "" {
			manufacturerMap[car.ManufacturerName] = true
		}

		if car.ManufacturerCountry != "" {
			countryMap[car.ManufacturerCountry] = true
		}

		if car.CategoryName != "" {
			categoryMap[car.CategoryName] = true
		}

		if car.Engine != "" {
			engineMap[car.Engine] = true
		}

		if car.Transmission != "" {
			transmissionMap[car.Transmission] = true
		}

		if car.Drivetrain != "" {
			drivetrainMap[car.Drivetrain] = true
		}
	}

	for manufacturer := range manufacturerMap {
		pageData.ManufacturerOptions = append(pageData.ManufacturerOptions, manufacturer)
	}

	for country := range countryMap {
		pageData.CountryOptions = append(pageData.CountryOptions, country)
	}

	for category := range categoryMap {
		pageData.CategoryOptions = append(pageData.CategoryOptions, category)
	}

	for engine := range engineMap {
		pageData.EngineOptions = append(pageData.EngineOptions, engine)
	}

	for transmission := range transmissionMap {
		pageData.TransmissionOptions = append(pageData.TransmissionOptions, transmission)
	}

	for drivetrain := range drivetrainMap {
		pageData.DrivetrainOptions = append(pageData.DrivetrainOptions, drivetrain)
	}

	return pageData
}

func filterMainCards(mainCards []Main, filters FilterOptions) []Main {
	var filtered []Main

	search := strings.ToLower(filters.Search)
	manufacturer := strings.ToLower(filters.Manufacturer)
	country := strings.ToLower(filters.Country)
	category := strings.ToLower(filters.Category)
	engine := strings.ToLower(filters.Engine)
	transmission := strings.ToLower(filters.Transmission)
	drivetrain := strings.ToLower(filters.Drivetrain)

	for _, car := range mainCards {
		if search != "" && !strings.Contains(strings.ToLower(car.Name), search) {
			continue
		}

		if manufacturer != "" && strings.ToLower(car.ManufacturerName) != manufacturer {
			continue
		}

		if country != "" && strings.ToLower(car.ManufacturerCountry) != country {
			continue
		}

		if category != "" && strings.ToLower(car.CategoryName) != category {
			continue
		}

		if engine != "" && strings.ToLower(car.Engine) != engine {
			continue
		}

		if transmission != "" && strings.ToLower(car.Transmission) != transmission {
			continue
		}

		if drivetrain != "" && strings.ToLower(car.Drivetrain) != drivetrain {
			continue
		}

		if filters.MinYear != "" {
			yearValue, err := strconv.Atoi(filters.MinYear)
			if err == nil && car.Year < yearValue {
				continue
			}
		}

		if filters.MinHorsepower != "" {
			horsepowerValue, err := strconv.Atoi(filters.MinHorsepower)
			if err == nil && car.Horsepower < horsepowerValue {
				continue
			}
		}

		filtered = append(filtered, car)
	}

	return filtered
}