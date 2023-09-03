package render

import (
	"embed"
	"fmt"
	"net/http"
	"strings"

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
	fmt.Fprint(w, string(file))
}

func DynamicRender(w http.ResponseWriter, r *http.Request, template string) {
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

	// Parse the template
	parsed_template := Parse(string(file))


	// Write the template to the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, s-maxage=86400")
	fmt.Fprint(w, parsed_template)
}

func Parse(file string) string {
	// VERY VERY Ruff first draft of the parser
	separation := strings.Split(file, "{{ --- }}")
	head := separation[0]
	content := separation[1]

	// Parse the head
	key := ""
	value := ""
	for _, line := range strings.Split(head, "\n") {
		if (line == "") {
			continue
		}
		if (strings.Contains(line, ":")) {
			key = strings.Split(line, ":")[0]
			value = strings.Split(line, ":")[1]
			content = strings.ReplaceAll(content, "{{ " + key + " }}", value)
		}
	}
	
	return content
}
