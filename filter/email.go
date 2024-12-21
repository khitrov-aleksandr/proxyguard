package filter

import (
	"regexp"
)

func BlockByEmail(email string) bool {
	return isGmail(email)
}

func isGmail(email string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@gmail.com$`, email)
	return matched
}
