package options

import (
	"fmt"

	ally_redis "github.com/tubstrr/ally/database/redis"
	ally_global "github.com/tubstrr/ally/global"
)

func GetAllyOptions() ally_global.AllyOptionsObject {
	AllyOptions := ally_global.AllyOptionsObject{
		SiteName: "",
		SiteUrl: "",
		Theme: "default",
	}

	// Check redis for the options
	redisSiteName, err := ally_redis.GetKey("ally-options-site-name")
	if (err != nil) {
		AllyOptions.SiteName = redisSiteName
	}
	redisSiteUrl, err := ally_redis.GetKey("ally-options-site-url")
	if (err != nil) {
		AllyOptions.SiteUrl = redisSiteUrl
	}
	redisTheme, err := ally_redis.GetKey("ally-options-theme")
	if (err != nil) {
		AllyOptions.Theme = redisTheme
	}

	// If they exist, return them
	if (AllyOptions.SiteName != "" && AllyOptions.SiteUrl != "" && AllyOptions.Theme != "") {
		return AllyOptions
	}

	// If they don't exist, check the database
	db := ally_global.Database

	siteOptions, e := db.Query(`
		SELECT * FROM ally_site_options
	`)

	if (e != nil) {
		panic(e)
	}

	defer siteOptions.Close()
	for siteOptions.Next() {
		var id int
		var name string
		var value string

		e = siteOptions.Scan(&id, &name, &value)
		if (e != nil) {
			panic(e)
		}

		if (name == "site_name") {
			AllyOptions.SiteName = value
			// Set the redis key
			ally_redis.SetKey("ally_options_site_name", value)
		}
		if (name == "site_url") {
			AllyOptions.SiteUrl = value
			// Set the redis key
			ally_redis.SetKey("ally_options_site_url", value)
		}
		if (name == "site_theme") {
			AllyOptions.Theme = value
			// Set the redis key
			ally_redis.SetKey("ally_options_site_theme", value)
		}
	}

	fmt.Println(AllyOptions)
	return AllyOptions
}