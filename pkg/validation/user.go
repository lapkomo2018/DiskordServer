package validation

import (
	"errors"
	"net/mail"
)

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	emailAddress, err := mail.ParseAddress(email)
	if !(err == nil && emailAddress.Address == email) {
		return errors.New("invalid email")
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}

	return nil
}
