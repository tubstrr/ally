package ally

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/tubstrr/ally/database"
	"github.com/tubstrr/ally/database/users"
	"github.com/tubstrr/ally/environment"
	"github.com/tubstrr/ally/network"
	"github.com/tubstrr/ally/render"
)

var Test = "test"

func Server() {
	// Check the environment variables
	environment.Check_environment()
	
	// Check the database
	database.CheckDatabase()
		
	// Get the environment variables for the server
	env := environment.Get_environment_variable("ALLY_ENVIRONMENT", "production")
	// Start the server	
	fmt.Println("Starting Ally server")
	if (env == "development") {
		port := environment.Get_environment_variable("ALLY_SERVER_PORT", "3000")
		fmt.Println("Running in development mode")
		Serve(port)
	} else if (env == "production") {
		port := environment.Get_environment_variable("ALLY_SERVER_PORT", "8080")
		fmt.Println("Running in production mode")
		Serve(port)
	} else {
		fmt.Println("Running in unknown mode")
		fmt.Println("Please set the environment variable ALLY_ENVIRONMENT to either 'development' or 'production'")
		return
	}
}

func Serve(port string) {
	fmt.Println("Server started on port " + port)

	// Handle all admin routes
	http.HandleFunc("/ally-admin", Admin)
	http.HandleFunc("/ally-admin/", Admin)
	// Admin template routes
	http.HandleFunc("/ally-admin/login", AdminLogin)
	http.HandleFunc("/ally-admin/create-account", AdminCreateAccount)
	// Admin form routes
	http.HandleFunc("/ally-admin/forms/auth", network.Authorization)	
	http.HandleFunc("/ally-admin/forms/create-account", network.CreateAccount)
	http.HandleFunc("/ally-admin/forms/logout", network.Logout)

	http.HandleFunc("/", Ally)
	
	http.ListenAndServe(":" + port, nil)
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

func Admin(w http.ResponseWriter, r *http.Request) {
	// If user is not logged in, redirect to login page
	loggedIn := network.IsUserLoggedIn(w,r)
	if (!loggedIn) {
		network.Redirect(w, r, "/ally-admin/login")
		return
	}

	render.DynamicRender(w, r, "/admin/pages/index.ally", false)
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	// If user is logged in, redirect to admin
	network.RedirectIfUserLoggedIn(w, r)

	render.DynamicRender(w, r, "/admin/pages/login.ally", false)
}

func AdminCreateAccount(w http.ResponseWriter, r *http.Request) {
	if (users.IsUserTableEmpty()) {
		render.DynamicRender(w, r, "/admin/pages/create-account.ally", false)
	} else {
		network.Redirect(w, r, "/ally-admin/login")
	}
}