package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"logger-service/cmd/api/data"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var dbClient *mongo.Client

func GetRouter(client *mongo.Client) *mux.Router {
	dbClient = client
	r := mux.NewRouter()
	r.HandleFunc("/insertLog", InsertLog)
	r.HandleFunc("/getAll", GetAllLogs)
	r.HandleFunc("/getRecordById/{id}", GetRecordById)
	r.HandleFunc("/updateRecordById", UpdateRecordById)
	r.HandleFunc("/dropCollection", DropCollection)

	return r

}

func InsertLog(w http.ResponseWriter, r *http.Request) {

	if dbClient == nil {
		fmt.Println("encountered error: ", errors.New("client is nil"))
	}

	var logEntry data.LogEntry
	err := json.NewDecoder(r.Body).Decode(&logEntry)
	if err != nil {
		w.Write([]byte("error decoding log entry"))
	}
	res := logEntry.Insert(logEntry, dbClient)

	w.Write([]byte(res))
}
func GetAllLogs(w http.ResponseWriter, r *http.Request) {
	var logEntry data.LogEntry
	logEntries, err := logEntry.GetAll(dbClient)
	if err != nil {
		w.Write([]byte("got error: " + err.Error()))
	}
	json.NewEncoder(w).Encode(logEntries)
}
func GetRecordById(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	var logEntry data.LogEntry
	logEntries, err := logEntry.GetEntryByID(values["id"], dbClient)
	if err != nil {
		w.Write([]byte("got error: " + err.Error()))
	}
	json.NewEncoder(w).Encode(logEntries)

}
func UpdateRecordById(w http.ResponseWriter, r *http.Request) {
	var logEntry data.LogEntry
	err := json.NewDecoder(r.Body).Decode(&logEntry)
	if err != nil {
		w.Write([]byte("error decoding log entry"))
	}
	err = logEntry.UpdateCollection(logEntry, dbClient)

	if err != nil {
		fmt.Println("error", err)
		w.Write([]byte("Got Error while updating log entry" + err.Error()))
	} else {
		w.Write([]byte(fmt.Sprintf("updated log entry %s", logEntry.Id)))
	}

}
func DropCollection(w http.ResponseWriter, r *http.Request) {
	var logEntry data.LogEntry
	err := logEntry.DropCollections(dbClient)
	if err != nil {
		w.Write([]byte("got error while dropping collection" + err.Error()))
	}
	w.Write([]byte("successfully dropped collection"))
}
