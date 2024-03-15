package api

import "regexp"

// checks that a given username is valid
func usernameIsValid(uName string) bool {
	// a username must have 3 < len < 30 and only contain alphanumerical characters
	is_alnum := regexp.MustCompile(`^[a-zA-Z0-9 ]+$`).MatchString(uName)
	if len(uName) < 3 || len(uName) > 30 || !is_alnum {
		return false
	}
	return true
}
