package routers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"logManager/logger"
	"logManager/models"
	"net/http"
)

var log = logger.Instance

func InitRouter(router *mux.Router) {
	router.HandleFunc("/server/{id}", GetOne).Methods("GET")
}

func GetOne(writer http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"]
	log.Print(id)
	model, err := models.GetServerConfigById(1)
	if err != nil {
		log.Print(err)
	}
	writer.WriteHeader(http.StatusOK)
	bytes, err := json.Marshal(model)
	writer.Write(bytes)
}
