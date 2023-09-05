package ally

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/tubstrr/ally/database"
	"github.com/tubstrr/ally/database/users"
	"github.com/tubstrr/ally/environment"
	ally_global "github.com/tubstrr/ally/global"
	"github.com/tubstrr/ally/global/lifecycle"
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
	// Set the global variables
	ally_global.W = w
	ally_global.R = r

	path := r.URL.Path
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	split := strings.Split(path, "/")
	if (split[0] == "ally-admin") {
		fmt.Println("You shouldn't be here")
		network.FourOhFour()
		return
	}
	
	defer ally_global.ResetGlobalVariables()
	render.HtmlRender("/front-end/index.html")
}

func Admin(w http.ResponseWriter, r *http.Request) {
	// Set the global variables
	lifecycle.Init(w, r)
	defer lifecycle.End()
	
	// If user is logged in, render admin
	// If not, redirect to login
	if (ally_global.ActiveUser.LoggedIn) {
		render.DynamicRender("/admin/pages/index.ally", false)
	} else {
		network.Redirect("/ally-admin/login")
	}
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	// Set the global variables
	lifecycle.Init(w, r)
	defer lifecycle.End()

	// If user is logged in, redirect to admin
	network.RedirectIfUserLoggedIn()

	render.AdminRender("/admin/pages/login.ally")
}

func AdminCreateAccount(w http.ResponseWriter, r *http.Request) {
	if (users.IsUserTableEmpty()) {
		defer ally_global.ResetGlobalVariables()
		render.DynamicRender("/admin/pages/create-account.ally", false)
	} else {
		network.Redirect("/ally-admin/login")
	}
}