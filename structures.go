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
	FoundingYear int `json:"foundingYear"`
}

type PageData struct {
	Filters FilterOptions

	ManufacturerOptions []string
	CountryOptions      []string
	CategoryOptions     []string
	EngineOptions       []string
	TransmissionOptions []string
	DrivetrainOptions   []string

	Mains []Main
	CompareCount int
}

type FilterOptions struct {
	Search        string
	Manufacturer  string
	Country       string
	Category      string
	MinYear       string
	Engine        string
	MinHorsepower string
	Transmission  string
	Drivetrain    string
}

type Main struct {
	ID             		 int
	Name           		 string
	Year           		 int
	Image          		 string
	ManufacturerName     string
	ManufacturerCountry  string
	ManufacturerYear  	 int
	CategoryName 		 string

	Engine       string
	Horsepower   int
	Transmission string
	Drivetrain   string
	IsCompared bool
}

type  dataResult struct {
	Name string
	Cars []Car
	Categories []Category
	Manufacturers []Manufacturer
	Error error
}

type ComparisonPageData struct {
	Cars []Main
}