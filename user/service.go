package user

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/zamedic/multiPoolDashboard/auth"
)

type Service interface {
	AddUser(email, password string) error
	ValidateUser(email, password string) (string, error)
	AddAPIEndpoint(email, pool, token string) error
}

type service struct {
	store Store
}

func NewService(store Store) Service {
	return &service{store: store}
}

func (s service) AddUser(email, password string) error {
	return s.store.addUser(email, password)
}

func (s service) ValidateUser(email, password string) (string, error) {
	err := s.store.findUser(email, password)
	if err != nil {
		return "", err
	}
	return generateToken(auth.Token,email)
}

func (s service) AddAPIEndpoint(email, pool, token string) error {
	return s.store.addKey(email,pool,token)
}

func generateToken(signingKey []byte, email string) (string, error) {
	claims := auth.CustomClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 120).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}
