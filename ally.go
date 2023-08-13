package ally

import (
		"fmt"
		"net/http"
		"strings"
		"github.com/tubstrr/ally/environment"
		"github.com/tubstrr/ally/render"
		"github.com/tubstrr/ally/network"
)

func Server() {
	// Check the environment variables
	environment.Check_environment()

	// Get the environment environment variable
	env := environment.Get_environment_variable("ALLY_ENVIRONMENT", "production")

	// Start the server	
	fmt.Println("Starting Ally server")
	if (env == "development") {
		fmt.Println("Running in development mode")
		Serve("3000")
		
	} else if (env == "production") {
		fmt.Println("Running in production mode")
		Serve("443")
	} else {
		fmt.Println("Running in unknown mode")
		fmt.Println("Please set the environment variable ALLY_ENVIRONMENT to either 'development' or 'production'")
		return
	}
}

func Serve(port string) {
	fmt.Println("Server started on port " + port)
	http.HandleFunc("/ally-admin", Admin)
	http.HandleFunc("/", Ally)
	http.ListenAndServe(":" + port, nil)
}
func Admin(w http.ResponseWriter, r *http.Request) {
	render.HtmlRender(w, r, "/admin/index.html")
}

func Ally(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	split := strings.Split(path, "/")
	if (split[0] == "ally-admin") {
		fmt.Println("You shouldn't be here")
		network.FourOhFour(w, r)
		return
	}

	render.HtmlRender(w, r, "/front-end/index.html")
}