package user

import (
	"github.com/go-kit/kit/endpoint"
	"context"
	"net/http"
	"encoding/json"
	"github.com/go-kit/kit/auth/jwt"

	jwt2 "github.com/dgrijalva/jwt-go"
)

func makeCreateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		u := request.(UserRequest)
		return nil, s.AddUser(u.Email, u.Password)
	}
}

func makeFindUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		u := request.(UserRequest)
		return s.ValidateUser(u.Email, u.Password)
	}
}

func makeAddApiEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		u := request.(ApiToken)
		claim := ctx.Value(jwt.JWTClaimsContextKey).(jwt2.MapClaims)
		return nil, s.AddAPIEndpoint(claim["email"].(string), u.Pool, u.Token)
	}
}

func decodeNewUserRequestJSON(_ context.Context, r *http.Request) (interface{}, error) {
	item := UserRequest{}
	err := json.NewDecoder(r.Body).Decode(&item)
	return item, err
}

func decodeApiTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	item := ApiToken{}
	err := json.NewDecoder(r.Body).Decode(&item)
	return item, err
}

type UserRequest struct {
	Email, Password string
}

type ApiToken struct {
	Pool, Token string
}
