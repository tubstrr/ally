package sessions

import (
	"time"

	"github.com/tubstrr/ally/database"
	"github.com/tubstrr/ally/database/users"
	"github.com/tubstrr/ally/errors"
)

func SetSessionToken(UserID int, SessionID string) {
	// Check if the user has a session token in the DB
	// if they do, update it
	// if they don't, create it
	db := database.OpenConnection()

	// Check if the user has a session token

	
	_, e := db.Exec(`
	INSERT INTO ally_user_sessions
			(user_id, session_id, created_at)
		VALUES
			($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
		SET
			session_id = excluded.session_id,
			created_at = excluded.created_at;
	`, UserID, SessionID, time.Now())
	errors.CheckError(e)


	
	defer database.CloseConnection(db)
}

func CheckSessionToken(session string) bool {
	// Check if the session token is valid
	db := database.OpenConnection()

	// Check if the user has a session token
	checkQuery, e := db.Query(`
		SELECT * FROM ally_user_sessions
		WHERE session_id = $1
	`, session)
	errors.CheckError(e)

	var id int
	var user_id int
	var session_id string
	var created_at string

	defer checkQuery.Close()
	for checkQuery.Next() {
		e = checkQuery.Scan(&id, &user_id, &session_id, &created_at)
		errors.CheckError(e)
	}

	defer database.CloseConnection(db)

	// if session is older than 10 seconds, delete it
	now := time.Now()
	// convert created_at to time
	converted_created_at, e := time.Parse(time.RFC3339, created_at)
	errors.CheckError(e)

	// TODO:
		// Find a way to have global options without using a billion
		// .env variables
	if (created_at != "" && now.Sub(converted_created_at).Hours() > 2) {
		DeleteSessionToken(session)
		return false
	}
	if (session_id == session) {
		return true
	} else {
		return false
	}
}

func GetUserIDFromSession(session string) int {
	// Check if the session token is valid
	db := database.OpenConnection()

	// Check if the user has a session token
	checkQuery, e := db.Query(`
		SELECT * FROM ally_user_sessions
		WHERE session_id = $1
	`, session)
	errors.CheckError(e)

	var id int
	var user_id int
	var session_id string
	var created_at string

	defer checkQuery.Close()
	for checkQuery.Next() {
		e = checkQuery.Scan(&id, &user_id, &session_id, &created_at)
		errors.CheckError(e)
	}

	defer database.CloseConnection(db)

	return user_id
}

func GetUserFromSession(session string) users.User {
	// Check if the session token is valid
	db := database.OpenConnection()

	// Check if the user has a session token
	checkQuery, e := db.Query(`
		SELECT * FROM ally_user_sessions
		WHERE session_id = $1
	`, session)
	errors.CheckError(e)

	var id int
	var user_id int
	var session_id string
	var created_at string

	defer checkQuery.Close()
	for checkQuery.Next() {
		e = checkQuery.Scan(&id, &user_id, &session_id, &created_at)
		errors.CheckError(e)
	}

	defer database.CloseConnection(db)

	return users.GetUserByID(user_id)
}

func DeleteSessionToken(session string) {
	// Delete the session token
	db := database.OpenConnection()

	// Check if the user has a session token
	_, e := db.Exec(`
		DELETE FROM ally_user_sessions
		WHERE session_id = $1
	`, session)
	errors.CheckError(e)

	defer database.CloseConnection(db)
}