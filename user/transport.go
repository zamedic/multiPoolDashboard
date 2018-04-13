package user

import (
	"github.com/go-kit/kit/log"
	"net/http"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/zamedic/multiPoolDashboard/gokit"
	gokit2 "github.com/weautomateeverything/go2hal/gokit"
	stdjwt "github.com/dgrijalva/jwt-go"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/zamedic/multiPoolDashboard/auth"
)

func MakeHandler(service Service, logger log.Logger) http.Handler {
	opts := gokit.GetServerOpts(logger)
	addUser := kithttp.NewServer(
		makeCreateUserEndpoint(service),
		decodeNewUserRequestJSON,
		gokit2.EncodeResponse,
		opts...
	)

	findUser := kithttp.NewServer(
		makeFindUserEndpoint(service),
		decodeNewUserRequestJSON,
		gokit2.EncodeResponse,
		opts...
	)

	apiKey := kithttp.NewServer(
		jwt.NewParser(auth.KeyFunction,stdjwt.SigningMethodHS256, jwt.MapClaimsFactory)(makeAddApiEndpoint(service)),
		decodeApiTokenRequest,
		gokit2.EncodeResponse,
		opts...
	)

	r := mux.NewRouter()
	r.Handle("/user/", addUser).Methods("POST")
	r.Handle("/user/auth",findUser).Methods("POST")
	r.Handle("/user/key",apiKey).Methods("POST")

	return r

}
