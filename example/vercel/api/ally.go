package handler

import (
	"fmt"
	"net/http"
)


func Ally(w http.ResponseWriter, r *http.Request) {
  theme := GetTheme()
  // Get the template
	file, err := theme.ReadFile("theme/templates/front-end/index.html")
  if (err != nil) {
    fmt.Println("Error reading file")
		fmt.Println(err)
		// network.FourOhFour(w, r)
		return
	}

	fmt.Println(string(file))

  w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, s-maxage=86400")
	fmt.Fprintf(w, string(file))

  // ally.Ally(w, r, theme)
}