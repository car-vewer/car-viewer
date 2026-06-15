package main
import (
        "fmt"
        "net/http"
        "io"
        "encoding/json"
        "strings"
        "strconv"
    )

type Specifications struct {
    Engine string `json:"engine"`
    Horsepower int `json:"horsepower"`
    Transmission string `json:"transmission"`
    Drivetrain string `json:"drivetrain"`
}

type Car struct {
    Id int `json:"id"`
    Name string `json:"name"`
    ManufacturerId int `json:"manufacturerId"`
    CategoryId int `json:"categoryId"`
    Year int `json:"year"`
    Specifications Specifications `json:"specifications"`
    Image string `json:"image"`
}

func carViewer(w http.ResponseWriter, r *http.Request) {
    response, err := http.Get("http://localhost:3000/api/models")
    if err != nil {
        fmt.Println("ERROR: ", err)
        return
    }
    carId := strings.TrimPrefix(r.URL.Path, "/cars/")
    id, err2 := strconv.Atoi(carId)
    if err2 != nil {
        fmt.Println("ERROR: ", err2)
        return
    }
    defer response.Body.Close()

    jsonData, err1 := io.ReadAll(response.Body)
        if err1 != nil {
        fmt.Println("ERROR: ", err1)
        return
    }
    var cars []Car
    json.Unmarshal(jsonData, &cars)
    for _, car := range cars {
        if car.Id == id {
            response.Header.Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(car)
            return
        }
    }
   // response.Header.Set("Content-Type", "application/json")
   // json.NewEncoder(w).Encode(car)
  //  w.Write([]byte(car.Name))
}

func manViewer(w http.ResponseWriter, r *http.Request) {
    response, err := http.Get("http://localhost:3000/api/manufacturers")
    if err != nil {
        fmt.Println("ERROR: ", err)
        return
    }
    defer response.Body.Close()

    jsonData, err1 := io.ReadAll(response.Body)
        if err1 != nil {
        fmt.Println("ERROR: ", err1)
        return
    }
    response.Header.Set("Content-Type", "application/json")
    w.Write(jsonData)
}

func catViewer(w http.ResponseWriter, r *http.Request) {
    response, err := http.Get("http://localhost:3000/api/categories")
    if err != nil {
        fmt.Println("ERROR: ", err)
        return
    }
    defer response.Body.Close()

    jsonData, err1 := io.ReadAll(response.Body)
        if err1 != nil {
        fmt.Println("ERROR: ", err1)
        return
    }
    response.Header.Set("Content-Type", "application/json")
    w.Write(jsonData)
}

func main() {
    http.HandleFunc("/cars/", carViewer)
    http.HandleFunc("/manufacturers", manViewer)
    http.HandleFunc("/categories", catViewer)
    render := "Server is running on http://localhost:8084/"
    fmt.Println(render)
    err := http.ListenAndServe(":8084", nil)
    if err != nil {
        fmt.Println("ERROR: ", err)
    }
}