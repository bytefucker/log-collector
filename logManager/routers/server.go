package routers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"logManager/models"
	"net/http"
)

func init() {
	Instance.HandleFunc("/server/{id}", getById).Methods("GET")
}

func getById(writer http.ResponseWriter, request *http.Request) {
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
