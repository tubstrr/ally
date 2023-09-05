package seo

import (
	ally_redis "github.com/tubstrr/ally/database/redis"
	ally_global "github.com/tubstrr/ally/global"
)

func GetTitle() string {
	title := ally_global.SEO.Title
	// Check Redis for the key "ally_option_site_name"
	if (title == "") {
		redisTitle, err := ally_redis.GetKey("ally_option_site_name")
		if (err != nil) {
			panic(err)
		}

		title = redisTitle
	}

	// If the title is still empty, check the database
	if (title == "") {
		// Get the site name from the database
		db := ally_global.Database

		siteNameQuery, e := db.Query(`
			SELECT option_value FROM ally_site_options
			WHERE option_name = 'site_name'
		`)
		if (e != nil) {
			panic(e)
		}

		defer siteNameQuery.Close()
		var siteName string
		for siteNameQuery.Next() {
			e = siteNameQuery.Scan(&siteName)
			if (e != nil) {
				panic(e)
			}
		}

		title = siteName
	}

	// If the title is still empty, set it to "Ally CMS"
	if (title == "") {
		title = "Ally CMS"
	}
	
	return title
}

func SetTitle(title string) {
	ally_global.SEO.Title = title
}

func SetDescription(description string) {
	ally_global.SEO.Description = description
}