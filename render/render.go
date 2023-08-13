package render

import (
		"fmt"
		"net/http"
		"os"
		"github.com/tubstrr/ally/network"
)

func StaticRender(w http.ResponseWriter, r *http.Request, template string) {
	file, err := os.ReadFile(template) // just pass the file name
  if err != nil {
		fmt.Print(err)
		network.FourOhFour(w, r)
		return
  }

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
  html := string(file) // convert content to a 'string'
  fmt.Fprintf(w, html) // write data to response
}