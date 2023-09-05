package lifecycle

import (
	"fmt"
	"net/http"

	"github.com/tubstrr/ally/database"
	"github.com/tubstrr/ally/database/options"
	ally_redis "github.com/tubstrr/ally/database/redis"
	"github.com/tubstrr/ally/database/seo"
	"github.com/tubstrr/ally/database/sessions"
	"github.com/tubstrr/ally/environment"
	ally_global "github.com/tubstrr/ally/global"
	"github.com/tubstrr/ally/network"
)

func Init(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Initializing Ally")
	
	// Reset the global variables
	ally_global.ResetGlobalVariables()

	// Set the global variables
	ally_global.W = w
	ally_global.R = r

	// Check the environment variables
	environment.Check_environment()
	
	// Check the redis database
	ally_redis.CheckRedisDatabase()
	
	// Check the database
	database.CheckDatabase()

	// Set global variables
	// -> Database
	ally_global.Database = database.OpenConnection()

	// -> User
	session := network.GetCookie("ally-user-session")
	user := sessions.GetUserFromSession(session)
	userLoggedIn := false
	if (user.Id != 0) {
		userLoggedIn = true
	}
	ally_global.ActiveUser = ally_global.ActiveUserObject{
		Session: session,
		LoggedIn: userLoggedIn,
		Id: user.Id,
		Username: user.Username,
		Email: user.Email,
		Role: user.Role,
	}

	// -> SEO
	ally_global.SEO = ally_global.SeoObject{
		Title: seo.GetTitle(),
		Description: "Ally CMS is a content management system built with Go",
	}

	// -> Ally Options
	ally_global.AllyOptions = options.GetAllyOptions()
}

func End() {
	ally_global.ResetGlobalVariables()
}