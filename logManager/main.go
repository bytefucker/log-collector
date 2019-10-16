package main

import (
	"github.com/gorilla/mux"
	_ "logManager/configs"
	"logManager/logger"
	"logManager/middlewares"
	"logManager/routers"
	"net/http"
	"time"
)

var log = logger.Instance

func main() {
	router := mux.NewRouter()

	router.Use(middlewares.LoggingMiddleware)
	router.Use(middlewares.ContentTypeMiddleware)

	routers.InitRouter(router)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Error(srv.ListenAndServe())
}
