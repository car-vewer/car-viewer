package main

type Specification struct {
	Engine string `json:"engine"`
	Horsepower int `json:"horsepower"`
	Transmission string `json:"transmission"`
	Drivetrain string `json:"drivetrain"`
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

type Category struct {
	ID int `json:"id"`
	Name string `json:"name"`
}

type Manufacturer struct {
	ID       int `json:"id"`
	Name string `json:"name"`
	Country  string `json:"country"`
}

type Main struct {
	ID             		 int
	Name           		 string
	Year           		 int
	Image          		 string
	ManufacturerName     string
	ManufacturerCountry  string
	CategoryName 		 string
}

type  dataResult struct {
	Name string
	Cars []Car
	Categories []Category
	Manufacturers []Manufacturer
	Error error
}