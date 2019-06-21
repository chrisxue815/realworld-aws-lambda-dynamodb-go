package model

import (
	"fmt"
	"os"
	"sync"
)

var once sync.Once
var userTableName string
var emailUserTableName string

func initializeSingletons() {
	prefix := fmt.Sprintf("realworld-%s", os.Getenv("STAGE"))
	userTableName = fmt.Sprintf("%s-user", prefix)
	emailUserTableName = fmt.Sprintf("%s-email-user", prefix)
}

func UserTableName() string {
	once.Do(initializeSingletons)
	return userTableName
}

func EmailUserTableName() string {
	once.Do(initializeSingletons)
	return emailUserTableName
}

type User struct {
	Username string
	Email    string
	Password []byte
	Image    string
	Bio      string
}

type EmailUser struct {
	Email    string
	Username string
}
