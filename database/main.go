package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // postgres driver

	"github.com/tubstrr/ally/environment"
	"github.com/tubstrr/ally/errors"
)

var db *sql.DB
var err error

func Check_database() {
	fmt.Println("Checking database")

	// Open a connection to the database
	db = OpenConnection()

	tables_needed := []string{
		"ally_users", 
		"ally_user_roles",
		"ally_user_sessions",
	}
	check_query := MakeCheckQuery(tables_needed)

	// Check if the database exists with the tables needed
	check, e := db.Query(check_query)
	errors.CheckError(e)

	defer check.Close()
	for check.Next() {
		var status bool

		err = check.Scan(&status)
		errors.CheckError(err)

		if (status) {
			fmt.Println("Database has all tables needed")
		} else {
			fmt.Println("Database does not have all tables needed")
			CreateDatabaseTables(tables_needed)
		}
	}

	// Close the connection
	defer CloseConnection(db)	
}

func OpenConnection() *sql.DB {
	connection_string := MakeConnectionString()
	db, err = sql.Open("postgres", connection_string)

	if err != nil { panic(err) }
	if err = db.Ping(); err != nil { panic(err) }

	fmt.Println("Database connection opened")
	return db
}

func CloseConnection(db *sql.DB) {
	db.Close()
	fmt.Println("Database connection closed")
}

// String functions
func MakeConnectionString() string {
	// Set up environment variables
	host := environment.Get_environment_variable("ALLY_DB_HOST", "localhost")
	port := environment.Get_environment_variable("ALLY_DB_PORT", "5432")
	db_name := environment.Get_environment_variable("ALLY_DB_NAME", "ally")
	user := environment.Get_environment_variable("ALLY_DB_USER", "ally")
	password := environment.Get_environment_variable("ALLY_DB_PASSWORD", "ally")

	// Run the check
	if (environment.Get_environment_variable("ALLY_ENVIRONMENT", "development") == "development") {
		db_name = db_name + "?sslmode=disable"
	}
	connection_string := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + db_name
	return connection_string
}

func MakeCheckQuery(tables []string) string {
	check_string := "SELECT EXISTS ("
	for i, table := range tables {
		if (i > 0) { check_string += " INTERSECT " }
		check_string += "(SELECT 1 FROM INFORMATION_SCHEMA.Tables WHERE TABLE_NAME='" + table + "')"
	}
	check_string += ");"

	return check_string
}

// Helper functions
func CreateDatabaseTables(tables []string) {
	fmt.Println("Creating database tables")
	// Define the tables schema
	tables_schema := map[string]string {
		"ally_users": `
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(50) NOT NULL,
			email VARCHAR(50) UNIQUE NOT NULL,
			role_id INT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		`,
		"ally_user_roles": `
			id SERIAL PRIMARY KEY,
			role VARCHAR(50) UNIQUE NOT NULL
		`,
		"ally_user_sessions": `
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			session_id VARCHAR(50) UNIQUE NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		`,
	}
	tables_preload := map[string]string {
		"ally_user_roles": `
			INSERT INTO ally_user_roles (role) VALUES ('superadmin') 
				ON CONFLICT DO NOTHING;
			INSERT INTO ally_user_roles (role) VALUES ('admin')
				ON CONFLICT DO NOTHING;
			INSERT INTO ally_user_roles (role) VALUES ('user')
				ON CONFLICT DO NOTHING;
		`,
	}

	// Loop through the tables and create them if they don't exist
	for _, table := range tables {
		fmt.Println("Creating table " + table)
		_, e := db.Exec("CREATE TABLE IF NOT EXISTS " + table + " (" + tables_schema[table] + ");")
		errors.CheckError(e)

		if (tables_preload[table] != "") {
			fmt.Println("Preloading table " + table)
			_, e := db.Exec(tables_preload[table])
			errors.CheckError(e)
		}
	}
}

func InsertRow(Table string, Columns string, Values string) {
	// Insert a row into the database
	fmt.Println("Inserting row into " + Table)
	_, e := db.Exec("insert into " + Table + " (" + Columns + ") values(" + Values + ")")
	errors.CheckError(e)
}

func DeleteRow(Table string, Columns string, Values string) {
	// Delete a row from the database
	fmt.Println("Deleting row from " + Table)
	_, e := db.Exec("DELETE FROM " + Table + " WHERE " + Columns + " = " + Values)
	errors.CheckError(e)
}
