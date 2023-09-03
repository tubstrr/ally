package ally_strings

import (
	"regexp"
	"strconv"
	"strings"
)


func KebabCase(input string) string {
	// Replace spaces and underscores with hyphens
	kebabString := strings.ReplaceAll(input, " ", "-")
	kebabString = strings.ReplaceAll(kebabString, "_", "-")

	// Convert to lowercase
	kebabString = strings.ToLower(kebabString)

	return kebabString
}

func AlphaNumeric(input string) string {
	// Define a regular expression to match non-alphanumeric characters
	regex := regexp.MustCompile("[^a-zA-Z0-9]+")

	// Replace non-alphanumeric characters with an empty string
	alphanumericString := regex.ReplaceAllString(input, "")

	return alphanumericString
}

func StringToNumber(input string) interface{} {
	// Try to convert the string to an integer first
	if intValue, err := strconv.Atoi(input); err == nil {
		return intValue
	}

	// If it's not an integer, try to convert it to a float
	if floatValue, err := strconv.ParseFloat(input, 64); err == nil {
		return floatValue
	}

	// Return an error if neither conversion succeeds
	return nil
}
