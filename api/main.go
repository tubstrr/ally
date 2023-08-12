package main


import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    static := http.FileServer(http.Dir("./static"))
    http.Handle("/", static)

    fmt.Printf("Starting server at port 443\n")
    if err := http.ListenAndServe(":443", nil); err != nil {
        log.Fatal(err)
    }
}