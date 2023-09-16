package main

import (
	"authentication-service/cmd/api/config"
	"authentication-service/cmd/api/router"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// const connectionString = "user=admin dbname=postgres sslmode=disable password=admin host=localhost port=5432"

func main() {
	//create class data/models.go
	//it will have user struct with name, password,status, created_at, modified_at password should not be visible, class should
	//function to add users, remove users and update user and validate userconst
	//create routes.go file to maintain routes
	//create mai
	connectionString := os.Getenv("DSN")
	connectionString = strings.Trim(connectionString, `"`)
	// connectionString := "user=admin dbname=postgres sslmode=disable password=admin host=127.0.0.1 port=5432"
	fmt.Println("DSN", connectionString)
	db := InitializeDB(connectionString)
	config.SetDb(db)
	r := router.GetMuxRouter(db)
	fmt.Println("got router")
	// http.ListenAndServe(":8080", r)
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening")

}
func InitializeDB(connectionString string) *sql.DB {

	Db, err := sql.Open("postgres", connectionString)
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)
	if err != nil {
		log.Fatal(err)
	}
	if Db != nil {
		fmt.Println("db connection created")
	}
	err = Db.Ping()
	fmt.Println("error while pinging", err)

	return Db
	// defer db.Close()
}
