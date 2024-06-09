package handlers

import "regexp"

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidateDate(date string) bool {
	return dateRegex.MatchString(date)
}
