package handler
 
import (
  "net/http"
  "github.com/tubstrr/ally"
)

func Ally(w http.ResponseWriter, r *http.Request) {
  ally.Ally(w, r)
}