package ally_global

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq" // postgres driver
)

var W http.ResponseWriter = nil
var R *http.Request = nil

// User information
type ActiveUserObject struct {
	Session string
	Id int
	LoggedIn bool
	Username string
	Email string
	Role int
}

var ActiveUser = ActiveUserObject{
	Session: "",
	Id: 0,
	LoggedIn: false,
	Username: "",
	Email: "",
	Role: 999,
}

// Database
var Database *sql.DB = nil

// Ally Options
type AllyOptionsObject struct {
	SiteName string
	SiteUrl string
	Theme string
}
var AllyOptions = AllyOptionsObject{
	SiteName: "",
	SiteUrl: "",
	Theme: "default",
}

// SEO
type SeoObject struct {
	Title string
	Description string
}
var SEO = SeoObject{}


func ResetGlobalVariables() {
	fmt.Println("Resetting global variables")
	W = nil
	R = nil

	// User information
	ActiveUser = ActiveUserObject{
		Session: "",
		Id: 0,
		LoggedIn: false,
		Username: "",
		Email: "",
		Role: 999,
	}

	// Database
	if (Database != nil) {
		Database.Close()
	}
	Database = nil

	// Ally Options
	AllyOptions = AllyOptionsObject{
		SiteName: "",
		SiteUrl: "",
		Theme: "default",
	}
}



