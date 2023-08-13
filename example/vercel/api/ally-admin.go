package handler
 
import (
  "net/http"
  "github.com/tubstrr/ally"
)

func AllyAdmin(w http.ResponseWriter, r *http.Request) {
  ally.Admin(w, r)
}