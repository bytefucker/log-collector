package routers

import (
	"github.com/gorilla/mux"
	"logManager/logger"
)

var log = logger.Instance

var Instance = mux.NewRouter().PathPrefix("/api").Subrouter()

func init() {
	Instance.Use(LoggingMiddleware)
	Instance.Use(ContentTypeMiddleware)
}
