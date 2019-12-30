package routers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"logManager/models"
	"net/http"
)

func init() {
	apiRouter.HandleFunc("/server/{id}", getById).Methods("GET")
	apiRouter.HandleFunc("/server", addServer).Methods("POST")
	apiRouter.HandleFunc("/server/list", queryList).Methods("POST")
}

func queryList(writer http.ResponseWriter, request *http.Request) {

}

func addServer(writer http.ResponseWriter, request *http.Request) {
	var server = models.ServerConfig{}
	params, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(params, &server)
	id, _ := models.AddServerConfig(&server)
	json.NewEncoder(writer).Encode(id)
}

func getById(writer http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"]
	model, err := models.GetServerConfigById(id)
	if err != nil {
		log.Print(err)
	}
	writer.WriteHeader(http.StatusOK)
	bytes, err := json.Marshal(model)
	writer.Write(bytes)
}
