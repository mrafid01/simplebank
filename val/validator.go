package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(str string, min, max int) error {
	n := len(str)
	if n > max || n < min {
		return fmt.Errorf("must contain from %d-%d characters", min, max)
	}
	return nil
}

func ValidateUsername(username string) error {
	if err := ValidateString(username, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(username) {
		return fmt.Errorf("must contain only lowercase letters, digits, or underscore")
	}
	return nil
}

func ValidatePassword(password string) error {
	return ValidateString(password, 6, 100)
}

func ValidateEmail(email string) error {
	if err := ValidateString(email, 3, 100); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("not valid email address")
	}
	return nil
}

func ValidateFullName(full_name string) error {
	if err := ValidateString(full_name, 3, 100); err != nil {
		return err
	}
	if !isValidFullName(full_name) {
		return fmt.Errorf("must contain only letters or spaces")
	}
	return nil
}
