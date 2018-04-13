package balance

import (
	"net/http"
	"github.com/go-kit/kit/log"
	"github.com/zamedic/multiPoolDashboard/gokit"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/zamedic/multiPoolDashboard/auth"
	http2 "github.com/go-kit/kit/transport/http"
	jwt2 "github.com/dgrijalva/jwt-go"
	gokit2 "github.com/weautomateeverything/go2hal/gokit"
	"github.com/gorilla/mux"
)

func MakeHandler(service Service, logger log.Logger) http.Handler {
	opts := gokit.GetServerOpts(logger)
	balanceEndpoint := http2.NewServer(
		jwt.NewParser(auth.KeyFunction,jwt2.SigningMethodHS256, jwt.MapClaimsFactory)(makeBalanceByPoolEndpoint(service)),
		gokit2.DecodeString,
		gokit2.EncodeResponse,
		opts...
	)

	coinBalanceEndpoint := http2.NewServer(
		jwt.NewParser(auth.KeyFunction,jwt2.SigningMethodHS256, jwt.MapClaimsFactory)(makeBalanceCoin(service)),
		decodeGetCoinBalanceRequest,
		gokit2.EncodeResponse,
		opts...
	)

	r := mux.NewRouter()
	r.Handle("/balance/pool", balanceEndpoint).Methods("GET")
	r.Handle("/balance/coin/{coin}",coinBalanceEndpoint).Methods("GET")

	return r
}
