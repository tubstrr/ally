package environment

import (
    "fmt"
    "os"
    "strings"
)

func Check_environment() {
	fmt.Println("Checking environment")
	// Set up checks
	keys_needed := [1]string {"ALLY_ENVIRONMENT"}
	keys_not_set := []string {}

	// Run the check
	keys_not_set = Loop_through_environment_variables(keys_needed, keys_not_set)
	
	// If there are any keys not set, load the .env file and check again
	if (len(keys_not_set) > 0) {
		Load_environment_variables()

		// Run the check again
		keys_not_set_after_load := []string {}
		keys_not_set_after_load = Loop_through_environment_variables(keys_needed, keys_not_set_after_load)

		// If there are still keys not set, print them out and exit
		if (len(keys_not_set_after_load) > 0) {
			fmt.Println("The following environment variables are not set:")
			for _, key := range keys_not_set_after_load {
				fmt.Println(" - " + key)
			}
			return
		}
	}
	fmt.Println("")

	return
}

func Loop_through_environment_variables(keys_needed [1]string, keys_not_set []string) []string {
	// Loop through the keys needed and check if they are set
	for _, key := range keys_needed {
		value := Get_environment_variable(key, "")
		// If the value is empty, add the key to the list of keys not set
		if (value == "") {
			keys_not_set = append(keys_not_set, key)
		}
	}
	return keys_not_set
}

func Get_environment_variable(key, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}

func Load_environment_variables() {
	fmt.Println("Loading environment variables")
	
	file, err := os.ReadFile(".env") // just pass the file name
  if err != nil {
		fmt.Print(err)
		return
  }

	text := string(file) // convert content to a 'string'
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		variable := strings.Split(line, "=")
		os.Setenv(variable[0], variable[1])
	}
}