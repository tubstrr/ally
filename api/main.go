package main

import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}

func main() {
	// Call the handler function in response to requests to the / path.
	http.HandleFunc("*", Handler)
}


// Vercel exmaple.
// package handler
 
// import (
//   "fmt"
//   "net/http"
// )
 
// func Handler(w http.ResponseWriter, r *http.Request) {
//   fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
// }