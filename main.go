package main
import (
        "fmt"
        "net/http"
        "io")

func carViewer(w http.ResponseWriter, r *http.Request) {
    response, err := http.Get("http://localhost:3000/api/models")
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
    http.HandleFunc("/", carViewer)
    render := "Server is running on http://localhost:8083"
    fmt.Println(render)
    err := http.ListenAndServe(":8083", nil)
    if err != nil {
        fmt.Println("ERROR: ", err)
    }
}