// package main
package handler
 
import (
  "fmt"
  "net/http"
)
 
func Handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}


// import (
//     "fmt"
//     "log"
//     "net/http"
// )

// func main() {
//     static := http.FileServer(http.Dir("./static"))
//     http.Handle("/", static)

//     fmt.Printf("Starting server at port 443\n")
//     if err := http.ListenAndServe(":443", nil); err != nil {
//         log.Fatal(err)
//     }
// }