package main

import (
	"github.com/zamedic/dynamodb"
	"net/http"
	"github.com/zamedic/multiPoolDashboard/user"
	"github.com/go-kit/kit/log"
	"os"
	"github.com/go-kit/kit/log/level"
	"os/signal"
	"syscall"
	"fmt"
	"github.com/zamedic/multiPoolDashboard/collector"
	"github.com/zamedic/multiPoolDashboard/balance"
)

func main() {
	db := dynamodb.NewConnection()

	userStore := user.NewDynamoStore(db)
	balanceStore := balance.NewDynamoStore(db)

	userService := user.NewService(userStore)
	balanceService := balance.NewService(balanceStore)

	collector.NewCollector(userStore, balanceStore)

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/user/", user.MakeHandler(userService, httpLogger))
	mux.Handle("/balance/",balance.MakeHandler(balanceService,httpLogger))
	http.Handle("/", accessControl(mux))

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", ":8000", "msg", "listening")
		errs <- http.ListenAndServe(":8000", nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
