package validators

import (
	"errors"
	"net/mail"
)

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("Invalid email")
	}
	return nil
}
