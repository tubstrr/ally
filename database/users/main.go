package users

import (
	_ "github.com/lib/pq" // postgres driver

	"github.com/tubstrr/ally/database"
	"github.com/tubstrr/ally/errors"
)

type User struct {
	Id int
	Role int
	Username string
	Email string
	Password string
}

func IsUserTableEmpty() bool {
	isTableEmpty := false

	// Check if the user exists
	db := database.OpenConnection()

	emptyQuery, e := db.Query(`SELECT CASE 
		WHEN EXISTS (SELECT * FROM ally_users LIMIT 1) THEN 1
		ELSE 0 
	END;`)
	errors.CheckError(e)

	var isEmpty int
	defer emptyQuery.Close()
	for emptyQuery.Next() {
		e = emptyQuery.Scan(&isEmpty)
		errors.CheckError(e)
		if (isEmpty == 0) {
			isTableEmpty = true
		}
	}

	defer database.CloseConnection(db)

	return isTableEmpty
}

func CreateUser(username string, email string, password string, role int) int {
	// Create the user
	db := database.OpenConnection()

	_, e := db.Exec(`
		INSERT INTO ally_users 
			(username, email, password, role_id) 
			VALUES 
				($1, $2, $3, $4)
			ON CONFLICT DO NOTHING;
	`, username, email, password, role)
	errors.CheckError(e)

	// Get the user
	userQuery, e := db.Query(`
		SELECT * FROM ally_users
		WHERE username = $1
	`, username)
	errors.CheckError(e)

	var id int
	var usernameQuery string
	var emailQuery string
	var passwordQuery string
	var role_idQuery int
	var created_atQuery string
	var updated_atQuery string

	defer userQuery.Close()
	for userQuery.Next() {
		e = userQuery.Scan(&id, &usernameQuery, &emailQuery, &passwordQuery, &role_idQuery, &created_atQuery, &updated_atQuery)
		errors.CheckError(e)
	}

	defer database.CloseConnection(db)

	return id
}

func IsValidUsername(username string) bool {
	isValid := false

	// Check if the user exists
	db := database.OpenConnection()

	userQuery, e := db.Query(`
		SELECT * FROM ally_users
		WHERE username = $1
	`, username)
	errors.CheckError(e)

	var id int
	var usernameQuery string
	var passwordQuery string
	var emailQuery string
	var role_idQuery int
	var created_atQuery string
	var updated_atQuery string

	defer userQuery.Close()
	for userQuery.Next() {
		e = userQuery.Scan(&id, &usernameQuery, &passwordQuery, &emailQuery, &role_idQuery, &created_atQuery, &updated_atQuery)
		errors.CheckError(e)
	}

	if (id == 0) {
		isValid = true
	}

	defer database.CloseConnection(db)

	return isValid
}

func IsValidEmail(email string) bool {
	isValid := false

	// Check if the user exists
	db := database.OpenConnection()

	userQuery, e := db.Query(`
		SELECT * FROM ally_users
		WHERE email = $1
	`, email)
	errors.CheckError(e)

	var id int
	var usernameQuery string
	var passwordQuery string
	var emailQuery string
	var role_idQuery int
	var created_atQuery string
	var updated_atQuery string

	defer userQuery.Close()
	for userQuery.Next() {
		e = userQuery.Scan(&id, &usernameQuery, &passwordQuery, &emailQuery, &role_idQuery, &created_atQuery, &updated_atQuery)
		errors.CheckError(e)
	}

	if (id == 0) {
		isValid = true
	}

	defer database.CloseConnection(db)

	return isValid
}

func GetUserByUsername(username string) User {
	// Get the user
	db := database.OpenConnection()

	userQuery, e := db.Query(`
		SELECT * FROM ally_users
		WHERE username = $1
	`, username)
	errors.CheckError(e)

	var id int
	var usernameQuery string
	var passwordQuery string
	var emailQuery string
	var role_idQuery int
	var created_atQuery string
	var updated_atQuery string

	defer userQuery.Close()
	for userQuery.Next() {
		e = userQuery.Scan(&id, &usernameQuery, &passwordQuery, &emailQuery, &role_idQuery, &created_atQuery, &updated_atQuery)
		errors.CheckError(e)
	}

	defer database.CloseConnection(db)

	user := User{
		Id: id,
		Username: usernameQuery,
		Email: emailQuery,
		Password: passwordQuery,
		Role: role_idQuery,
	}

	return user
}

func GetUserByID(id int) User {
	// Get the user
	db := database.OpenConnection()

	userQuery, e := db.Query(`
		SELECT * FROM ally_users
		WHERE id = $1
	`, id)
	errors.CheckError(e)

	var idQuery int
	var usernameQuery string
	var passwordQuery string
	var emailQuery string
	var role_idQuery int
	var created_atQuery string
	var updated_atQuery string

	defer userQuery.Close()
	for userQuery.Next() {
		e = userQuery.Scan(&idQuery, &usernameQuery, &passwordQuery, &emailQuery, &role_idQuery, &created_atQuery, &updated_atQuery)
		errors.CheckError(e)
	}

	defer database.CloseConnection(db)

	user := User{
		Id: idQuery,
		Username: usernameQuery,
		Email: emailQuery,
		Role: role_idQuery,
	}


	return user
}