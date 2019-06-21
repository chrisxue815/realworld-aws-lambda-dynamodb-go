package service

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"
	"strings"
	"time"
)

const tokenExpirationDay = 60
const passwordKeyLength = 64

var passwordSalt = []byte("KU2YVXA7BSNExJIvemcdz61eL86IJDCC")
var jwtSecret = []byte("C92cw5od80NCWIvu4NZ8AKp5NyTbnBmG")

func Scrypt(password string) ([]byte, error) {
	// https://godoc.org/golang.org/x/crypto/scrypt
	key, err := scrypt.Key([]byte(password), passwordSalt, 32768, 8, 1, passwordKeyLength)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func GenerateToken(username string) (string, error) {
	now := time.Now().UTC()
	exp := now.AddDate(0, 0, tokenExpirationDay).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": exp,
	})

	return token.SignedString(jwtSecret)
}

func VerifyAuthorization(auth string) (string, string, error) {
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Token" {
		return "", "", errors.New("invalid authorization")
	}

	token := parts[1]

	username, err := VerifyToken(token)
	return username, token, err
}

func VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, validateToken)

	if err != nil {
		return "", err
	}

	if token == nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token: invalid claims")
	}

	if !claims.VerifyExpiresAt(time.Now().UTC().Unix(), true) {
		return "", errors.New("invalid token: token expired")
	}

	username, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid token: sub missing")
	}

	return username, nil
}

func validateToken(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return jwtSecret, nil
}
