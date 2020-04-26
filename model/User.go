package model

import (
	"fmt"
)

const MinPasswordLength = 0
const PasswordKeyLength = 64

type User struct {
	Username     string
	Email        string
	PasswordHash []byte
	Image        string
	Bio          string
}

type EmailUser struct {
	Email    string
	Username string
}

type Follow struct {
	Follower  string
	Publisher string
}

func (u *User) Validate() error {
	if u.Username == "" {
		return NewInputError("username", "can't be blank")
	}

	if u.Email == "" {
		return NewInputError("email", "can't be blank")
	}

	if u.PasswordHash == nil || len(u.PasswordHash) != PasswordKeyLength {
		return NewInputError("password", "can't be blank")
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return NewInputError("password", fmt.Sprintf("must be at least %d characters in length", MinPasswordLength))
	}

	return nil
}
