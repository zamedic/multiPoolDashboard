package auth

import (
	"github.com/dgrijalva/jwt-go"
	stdjwt "github.com/dgrijalva/jwt-go"

)

type CustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var Token = []byte("supersecret")
var KeyFunction = func(token *stdjwt.Token) (interface{}, error) { return Token, nil }
