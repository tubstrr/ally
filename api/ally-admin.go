package handler
 
import (
  "fmt"
  "net/http"
  "os"
)

 
func AllyAdmin(w http.ResponseWriter, r *http.Request) {
  // fmt.Fprintf(w, "<h1>ally-admin.go</h1>") // write data to response

  backtrack := "../../../../"



  // currentDir, err := os.Getwd()
  // if err != nil {
  //   fmt.Println(err)
  // }
  // fmt.Println(currentDir) // for example /home/user

  filepath := backtrack + "ally/admin/index.html"

  file, err := os.ReadFile(filepath) // just pass the file name
  if err != nil {
      fmt.Print(err)
  }
  // fmt.Println(file) // print the content as 'bytes'
  str := string(file) // convert content to a 'string'
  fmt.Fprintf(w, str) // write data to response

  // fmt.Println(str) // print the content as a 'string'

  // Directory output helper
  // dir, err := os.Open(filepath)
  //  if err != nil {
  //     fmt.Println("Error:", err)
  //     return
  //  }
  //  defer dir.Close()
  //  files, err := dir.Readdir(-1)
  //  if err != nil {
  //     fmt.Println("Error:", err)
  //     return
  //  }
  //  for _, file := range files {
  //     fmt.Println(file.Name())
  //  }
  
}