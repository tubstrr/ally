package render

import (
	"embed"
	"fmt"
	"net/http"
	"strings"

	ally_global "github.com/tubstrr/ally/global"
	"github.com/tubstrr/ally/network"
)

var (
	//go:embed templates
	templates embed.FS
)

func HtmlRender(template string) {
	w := ally_global.W
	// r := ally_global.R
	
	// Check if the template exists
	if (template == "") {
		network.FourOhFour()
		return
	}

	// Get the template
	file, err := templates.ReadFile("templates" + template)
	if (err != nil) {
		fmt.Println(err)
		network.FourOhFour()
		return
	}

	// Write the template to the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, s-maxage=86400")
	if (network.IsUserLoggedIn()) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(file))
}

func DynamicRender(template string, cache bool) {
	w := ally_global.W 
	// r := ally_global.R 
	// Check if the template exists
	if (template == "") {
		network.FourOhFour()
		return
	}

	// Get the template
	file, err := templates.ReadFile("templates" + template)
	if (err != nil) {
		fmt.Println(err)
		network.FourOhFour()
		return
	}

	// Parse the template
	parsed_template := Parse(string(file))

	// Write the template to the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, s-maxage=86400")

	if (network.IsUserLoggedIn() || !cache) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, parsed_template)
}

func Parse(file string) string {
	// VERY VERY Ruff first draft of the parser
	// Get the template
	root, err := templates.ReadFile("templates/admin/root.ally")
	if (err != nil) {
		fmt.Println(err)
	}

	// Split the file into the head and the content
	content := file
	if (strings.Contains(file, "{{ --- }}")) {
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
	}

	// Inject the content into the root
	content = strings.ReplaceAll(string(root), "{{ ALLY_PAGE }}", content)

	if (ally_global.ActiveUser.LoggedIn) {
		content = strings.ReplaceAll(content, "{{ ALLY_USERNAME }}", ally_global.ActiveUser.Username)
		content = strings.ReplaceAll(content, "{{ ALLY_EMAIL }}", ally_global.ActiveUser.Email)
	}

	// Inject the SEO into the content
	SEO := ally_global.SEO

	content = strings.ReplaceAll(content, "{{ ALLY_SEO_TITLE }}", SEO.Title)
	content = strings.ReplaceAll(content, "{{ ALLY_SEO_DESCRIPTION }}", SEO.Description)

	return content
}

func AdminRender(template string) {
	w := ally_global.W 

	// Get the template
	file, err := templates.ReadFile("templates" + template)
	if (err != nil) {
		fmt.Println(err)
		network.FourOhFour()
		return
	}

	// Parse the template
	parsed_template := AdminParser(string(file))

	// Write the template to the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, must-revalidate, max-age=0, no-store")
	w.WriteHeader(http.StatusOK)
	
	// Actually render the template
	fmt.Fprint(w, parsed_template)
}

func AdminParser(file string) string {
	// VERY VERY Ruff first draft of the parser
	// Get the template
	root, err := templates.ReadFile("templates/admin/root.ally")
	if (err != nil) {
		fmt.Println(err)
	}

	// Split the file into the head and the content
	content := file
	if (strings.Contains(file, "{{ --- }}")) {
		separation := strings.Split(file, "{{ --- }}")
		head := separation[0]
		content = separation[1]
	
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
				// If value starts with "ally_global", convert it from string to the global variable
				content = strings.ReplaceAll(content, "{{ " + key + " }}", value)
			}
		}
	}

	// Inject the content into the root
	content = strings.ReplaceAll(string(root), "{{ ALLY_PAGE }}", content)

	islandsMap := map[string]string{
		"ALLY_ACTION_BAR": "templates/admin/islands/action-bar.ally",
	}

	for key, value := range islandsMap {
		island, err := templates.ReadFile(value)
		if (err != nil) {
			fmt.Println(err)
		}
		fmt.Println(key)
		fmt.Println(string(island))
		content = strings.ReplaceAll(content, "{{ " + key + " }}", string(island))
	}


	if (ally_global.ActiveUser.LoggedIn) {
		content = strings.ReplaceAll(content, "{{ ALLY_USERNAME }}", ally_global.ActiveUser.Username)
		content = strings.ReplaceAll(content, "{{ ALLY_EMAIL }}", ally_global.ActiveUser.Email)
	}

	// Inject the SEO into the content
	SEO := ally_global.SEO

	content = strings.ReplaceAll(content, "{{ ALLY_SEO_TITLE }}", SEO.Title)
	content = strings.ReplaceAll(content, "{{ ALLY_SEO_DESCRIPTION }}", SEO.Description)

	// Inject the Ally Options into the content
	content = strings.ReplaceAll(content, "{{ ALLY_SITE_NAME }}", ally_global.AllyOptions.SiteName)
	content = strings.ReplaceAll(content, "{{ ALLY_SITE_URL }}", ally_global.AllyOptions.SiteUrl)
	content = strings.ReplaceAll(content, "{{ ALLY_THEME }}", ally_global.AllyOptions.Theme)

	return content
}