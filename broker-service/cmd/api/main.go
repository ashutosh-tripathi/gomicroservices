package main

import (
	"fmt"
	"net/http"

	// Import the correct package
	"broker-service/cmd/api/router"

	"github.com/go-chi/cors"
)

func main() {
	// r := GetMuxRouter() // Call the exported function
	r := router.GetMuxRouter()
	fmt.Println("Starting go server...")
	handler := cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST"},
	})

	http.Handle("/", handler(r))
	http.ListenAndServe(":8080", nil)
}
