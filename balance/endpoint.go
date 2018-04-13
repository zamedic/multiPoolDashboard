package balance

import (
	"github.com/go-kit/kit/endpoint"
	"context"
	"github.com/go-kit/kit/auth/jwt"
	jwt2 "github.com/dgrijalva/jwt-go"
	"net/http"
	"github.com/gorilla/mux"
)

func makeBalanceByPoolEndpoint(s Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		claim := ctx.Value(jwt.JWTClaimsContextKey).(jwt2.MapClaims)
		r := s.getBalanceByPool(claim["email"].(string))
		return r,nil
	}
}

func makeBalanceCoin(s Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		coin := request.(string)
		claim := ctx.Value(jwt.JWTClaimsContextKey).(jwt2.MapClaims)
		return s.getBalanceByCoin(claim["email"].(string),coin), nil
	}
}

func decodeGetCoinBalanceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return mux.Vars(r)["coin"], nil
}