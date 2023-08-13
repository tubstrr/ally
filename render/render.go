package render

import (
		"fmt"
		"net/http"
		"embed"
		"github.com/tubstrr/ally/network"
)

var (
	//go:embed templates
	templates embed.FS
)

func HtmlRender(w http.ResponseWriter, r *http.Request, template string) {
	// Check if the template exists
	if (template == "") {
		network.FourOhFour(w, r)
		return
	}

	// Get the template
	file, err := templates.ReadFile("templates" + template)
	if (err != nil) {
		fmt.Println(err)
		network.FourOhFour(w, r)
		return
	}

	// Write the template to the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, s-maxage=86400")
	fmt.Fprintf(w, string(file))
}