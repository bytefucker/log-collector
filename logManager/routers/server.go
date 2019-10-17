package routers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"logManager/models"
	"net/http"
	"strconv"
)

func init() {
	apiRouter.HandleFunc("/server/{id}", getById).Methods("GET")
}

func getById(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(request)["id"])
	model, err := models.GetServerConfigById(id)
	if err != nil {
		log.Print(err)
	}
	writer.WriteHeader(http.StatusOK)
	bytes, err := json.Marshal(model)
	writer.Write(bytes)
}
