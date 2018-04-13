package gokit

import (
	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/kit/log"
	gokitjwt "github.com/go-kit/kit/auth/jwt"

	"github.com/weautomateeverything/go2hal/gokit"
)

func GetServerOpts(logger kitlog.Logger) []kithttp.ServerOption {
	return []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(gokit.EncodeError),
		kithttp.ServerBefore(gokitjwt.HTTPToContext()),

	}
}
