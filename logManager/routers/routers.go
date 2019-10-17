package routers

import (
	"github.com/gorilla/mux"
	"logManager/logger"
	"net/http"
)

var Router = mux.NewRouter()
var (
	log       = logger.Instance
	apiRouter = Router.PathPrefix("/api").Subrouter()
)

func init() {
	Router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("webapp/dist"))))
	Router.Use(LoggingMiddleware)
	apiRouter.Use(ContentTypeMiddleware)
}
