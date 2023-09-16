package router

import (
	"authentication-service/cmd/api/config"
	"authentication-service/cmd/api/data"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func GetMuxRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	fmt.Println("create router")
	r.HandleFunc("/createUser", createUserHandler)
	r.HandleFunc("/getUser/{id}", fetchUserHandler)
	r.HandleFunc("/getAllUser", fetchAllUserHandler)
	r.HandleFunc("/updateUser/{id}", updateUserHandler)
	r.HandleFunc("/deleteUser/{id}", deleteUser)
	r.HandleFunc("/userExists/{id}/{passwd}", checkUserExists)
	return r

}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user data.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("got error decoding", err)
	}
	fmt.Println("decoded user", user)
	user.Created_at = time.Now()
	user.Modified_at = time.Now()
	db := config.GetDb()
	if db == nil {
		fmt.Println("db still nill")
	}
	userid, err := user.AddUser(db)
	if err != nil {
		fmt.Println("Got error creating user", err)
	} else {
		fmt.Println("created user", userid)
	}
	json.NewEncoder(w).Encode(&user)
}
func fetchUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if val, ok := vars["id"]; ok {
		var user data.User
		fmt.Println("passed key is ", val)
		data, _ := user.GetUser(val, config.GetDb())
		json.NewEncoder(w).Encode(data)
	} else {
		fmt.Println("Please provide id")
	}

}
func fetchAllUserHandler(w http.ResponseWriter, r *http.Request) {
	var user data.User
	var users []data.User
	users, err := user.GetAll(config.GetDb())
	if err != nil {
		fmt.Println("got err: ", err)
	}
	json.NewEncoder(w).Encode(users)

}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user, updatedUser data.User
	vals := mux.Vars(r)
	user.UserID = vals["id"]
	fmt.Println("fetched id", user.UserID)
	json.NewDecoder(r.Body).Decode(&updatedUser)
	updatedUser.Created_at = time.Now()
	updatedUser.Modified_at = time.Now()
	fmt.Printf("updated user %+v", updatedUser)
	upd, err := user.UpdateUser(updatedUser, config.Db)
	if err != nil {
		fmt.Println("err", err)
	} else {

		w.Write([]byte(fmt.Sprintf("updated %d rows", upd)))
	}

}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	var user data.User
	vals := mux.Vars(r)
	user.UserID = vals["id"]
	user.DeleteUser(config.GetDb())
}
func checkUserExists(w http.ResponseWriter, r *http.Request) {
	var user data.User
	vals := mux.Vars(r)
	user.UserID = vals["id"]
	fmt.Println("user exists call with ", vals["id"], vals["passwd"])
	exists, _ := user.UserExists(config.GetDb(), vals["id"], vals["passwd"])
	if exists {
		logAuthenticated(vals["id"], vals["passwd"])
	}
	user.Status = exists
	json.NewEncoder(w).Encode(user)
	// w.Write([]byte(fmt.Sprintf("The user  %v", exists)))

}
func logAuthenticated(id, passwd string) {

	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = "auth_entry"
	entry.Data = "authenticated" + id + " with password " + passwd
	resbytes, _ := json.Marshal(entry)
	resp, err := http.Post("http://dockompose-logger-service-1/insertLog", "application/json", bytes.NewReader(resbytes))
	if err != nil {
		fmt.Println("got error in handleLogRequest:", err)

	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK status code:", resp.Status)
		// Handle non-OK status codes appropriately
		// Example: Return an error response

		return // Exit the function
	}
	defer resp.Body.Close()

}
