package main

import (
	"logManager/logger"
	"logManager/routers"
	"net/http"
	"time"
)

var log = logger.Instance

var router = routers.Instance

func main() {
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Error(srv.ListenAndServe())
}
