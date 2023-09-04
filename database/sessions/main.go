package sessions

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tubstrr/ally/database"
	ally_redis "github.com/tubstrr/ally/database/redis"
	"github.com/tubstrr/ally/database/users"
	"github.com/tubstrr/ally/errors"
)

type Session struct {
	ID int `json:"ID"`
	UserID int `json:"UserID"`
	SessionID string `json:"SessionID"`
	CreatedAt string `json:"CreatedAt"`
}

func SetSessionToken(UserID int, SessionID string) {
	// Check if the user has a session token in the DB

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

	// JSON encode the session
	sessionObject := Session{
		UserID: UserID,
		SessionID: SessionID,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	sessionJSON, _ := json.Marshal(sessionObject)

	// Set the session token in redis
	ally_redis.SetKey("session-" + SessionID, string(sessionJSON))
	
	defer database.CloseConnection(db)
}

func CheckSessionToken(session string) bool {
	// Variables used in this function
	now := time.Now()
	validSession := false
	
	// Check if the session token is valid
	redisSession, err := ally_redis.GetKey("session-" + session)
	if (err != nil) {
		validSession = false
	}

	if (redisSession != "") {
		var redisSessionObject Session
		err := json.Unmarshal([]byte(redisSession), &redisSessionObject)
		if (err != nil) {
			validSession = false
		}

		sessionIDsMatch := redisSessionObject.SessionID == session
		converted_created_at, e := time.Parse(time.RFC3339, redisSessionObject.CreatedAt)
		if (e != nil) {
			validSession = false
		}

		sessionIsExpired := redisSessionObject.CreatedAt != "" && now.Sub(converted_created_at).Hours() > 2
		if (sessionIDsMatch && !sessionIsExpired) {
			validSession = true
		} else {
			validSession = false
			DeleteSessionToken(session)
			ally_redis.DeleteKey("session-" + session)
		}
		return validSession
	}

	// Now we know the session is NOT in redis,
	// so we need to check the database

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
	// convert created_at to time
	if (created_at == "") {
		validSession = false
		// DeleteSessionToken(session)
		// ally_redis.DeleteKey("session-" + session)
		return validSession
	}

	converted_created_at, e := time.Parse(time.RFC3339, created_at)
	errors.CheckError(e)

	// TODO:
	// Find a way to have global options without using a billion
	// .env variables
	if (now.Sub(converted_created_at).Hours() > 2) {
		fmt.Print("Session is expired")
		// DeleteSessionToken(session)
		// ally_redis.DeleteKey("session-" + session)
	}
	if (session_id == session) {
		validSession = true
		// JSON encode the session
		sessionObject := Session{
			ID: id,
			UserID: user_id,
			SessionID: session_id,
			CreatedAt: created_at,
		}
		sessionJSON, _ := json.Marshal(sessionObject)
		// Set the session token in redis
		ally_redis.SetKey("session-" + session, string(sessionJSON))
	}

	return validSession
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
	// First check if the user is in redis
	redisSession, err := ally_redis.GetKey("session-" + session)
	if (err != nil) {
		fmt.Println(err)
	}

	if (redisSession != "") {
		var redisSessionObject Session
		err := json.Unmarshal([]byte(redisSession), &redisSessionObject)
		if (err != nil) {
			fmt.Println(err)
		}

		user := users.GetUserByID(redisSessionObject.UserID)
		return user
	}

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