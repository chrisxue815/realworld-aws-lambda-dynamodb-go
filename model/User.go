package model

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

type Follow struct {
	Follower string
	Publisher string
}
