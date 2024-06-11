package email

import (
	"errors"
	"net/mail"
)

func Validate(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email")
	}
	return nil
}
