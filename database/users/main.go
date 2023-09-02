package users

import (
	"fmt"

	"github.com/tubstrr/ally/database"
)

type User struct {
	Id int
	Username string
	Password string
	Email string
	Role string
}

func Create_ally_user(User) {
	fmt.Println("Creating ally_user")

	// Check if the user exists
	db := database.Open_connection()

	// If the user exists, update the user
	// If the user does not exist, create the user
	
	defer database.Close_connection(db)
}
