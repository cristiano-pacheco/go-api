package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// NotEmpty Implement a Required method to check that specific fields in the form
// data are present and not blank. If any fields fail this check, add the
// appropriate message to the form errors.
func NotEmpty(field string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s cannot be empty", field)
	}
	return nil
}

// MaxLength Implement a MaxLength method to check that a specific field in the form
// contains a maximum number of characters. If the check fails then add the
// appropriate message to the form errors.
func MaxLength(field string, value string, d int) error {
	if utf8.RuneCountInString(value) > d {
		return fmt.Errorf("%s is too long (maximum is %d characters)", field, d)
	}
	return nil
}

// MinLength Implement a MinLength method to check that a specific field in the form
// contains a minimum number of characters. If the check fails then add the
// appropriate message to the form errors.
func MinLength(field string, value string, d int) error {
	if utf8.RuneCountInString(value) < d {
		return fmt.Errorf("%s is too short (minimum is %d characters)", field, d)
	}
	return nil
}

// EmailRX regex rule
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Email validate email
func Email(field string, value string) error {
	if IsEmail(value) == false {
		return fmt.Errorf("%s is not a valid email", field)
	}
	return nil
}

// IsEmail validate the string email pattern
func IsEmail(value string) bool {
	err := MatchesPattern(value, emailRegex)
	if err != nil {
		return false
	}
	return true
}

// MatchesPattern Implement a MatchesPattern method to check that a specific field in the form
// matches a regular expression. If the check fails then add the
// appropriate message to the form errors.
func MatchesPattern(value string, pattern *regexp.Regexp) error {
	if !pattern.MatchString(value) {
		return fmt.Errorf("is invalid")
	}
	return nil
}
