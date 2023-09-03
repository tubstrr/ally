package validation

import (
	"net/mail"

	errors "github.com/tubstrr/ally/errors"
	ally_strings "github.com/tubstrr/ally/utilities/strings"
	"golang.org/x/crypto/bcrypt"
)

func ValidateEmail(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}

func ConvertUsername(username string) string {
	username = ally_strings.AlphaNumeric(username)
	username = ally_strings.KebabCase(username)
	return username
}

func ConvertPassword(password string) string {
	password = ally_strings.AlphaNumeric(password)
	password, e := HashPassword(password)
	errors.CheckError(e)
	// TODO: Add more password validation


	return password
}

func HashPassword(password string) (string, error) {
	// Generate a random salt with a cost factor of 12 (adjust the cost according to your needs)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword, inputPassword string) bool {
	// Compare the input password with the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}