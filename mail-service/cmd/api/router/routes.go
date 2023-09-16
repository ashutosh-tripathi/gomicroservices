package router

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mail-service/cmd/api/mailer"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	SMTPHost string
	SMTPPort int
}

var config Config

func GetRouter(_SMTPHost string, _SMTPPort int) *mux.Router {
	r := mux.NewRouter()
	r.Use(securityMiddleware)
	r.HandleFunc("/sendMail", SendMail)

	config = Config{
		SMTPHost: _SMTPHost,
		SMTPPort: _SMTPPort,
	}
	return r

}
func securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. CORS Handling
		w.Header().Set("Access-Control-Allow-Origin", "https://*,http://*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 2. CSRF Protection (Example: Generate and validate CSRF tokens)
		// Generate a CSRF token and store it in a cookie
		csrfToken := generateCSRFToken()
		// fmt.Println("CSRF token generated" + csrfToken)
		http.SetCookie(w, &http.Cookie{
			Name:  "csrf_token",
			Value: csrfToken,
			// Set secure and HttpOnly flags as needed
		})

		// Verify CSRF token in incoming requests
		// clientCSRFToken := r.Header.Get("X-CSRF-Token")
		// // if clientCSRFToken != csrfToken {
		// // 	http.Error(w, "CSRF token validation failed", http.StatusForbidden)
		// // 	return
		// // }

		// 3. Content Security Policy (CSP)
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' trusted-scripts.com")

		// 4. HTTP Strict Transport Security (HSTS)
		// w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// 5. X-Content-Type-Options
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// 6. X-Frame-Options
		w.Header().Set("X-Frame-Options", "DENY")

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func generateCSRFToken() string {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(tokenBytes)
	// return "csrftoken"

}

func SendMail(w http.ResponseWriter, r *http.Request) {
	var msg mailer.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	fmt.Println("in mail got response: ", msg)
	if err != nil {
		fmt.Println("Got error decoding")
		w.Write([]byte(err.Error()))
		return
	}
	var msgServ mailer.MessageServer
	msgServ.SMTPHOST = config.SMTPHost
	msgServ.SMTPPORT = config.SMTPPort
	err = msgServ.SendMessage(&msg)
	if err != nil {
		fmt.Println("Got error:", err)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("successfully sent message"))

}
