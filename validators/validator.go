package validators

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidString(value string, minLength int, maxLength int) error {
	if len(value) < minLength || len(value) > maxLength {
		return fmt.Errorf("string length must be between %d and %d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidString(value, 3, 20); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("username can only contain lower case letters, numbers, and underscores")
	}
	return nil
}

func ValidateFullname(value string) error {
	if err := ValidString(value, 3, 20); err != nil {
		return err
	}
	if !isValidFullname(value) {
		return fmt.Errorf("fullname can only contain letters and spaces")
	}
	return nil
}

func ValidateEmail(value string) error {
	if err := ValidString(value, 5, 200); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("invalid email address")
	}

	return nil
}

func ValidatePassword(value string) error {
	if err := ValidString(value, 6, 100); err != nil {
		return fmt.Errorf("password must be between 6 and 100 characters")
	}
	return nil
}
