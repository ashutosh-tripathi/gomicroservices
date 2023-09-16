package router

import (
	"broker-service/cmd/api/grpcproto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/rpc"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//create a common struct here, that willbe common for all requests, make action and authpayload, call auth service based on action

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Body    string `json:"body"`
}

type BrokerRequest struct {
	Action      string      `json:"action"`
	Payload     AuthPayload `json:"payload"`
	LogPayload  LogPayload  `json:"logPayload"`
	MailPayload MailPayload `json:"mailPayload"`
}
type AuthPayload struct {
	UserID      string    `json:"user"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	Status      bool      `json:"status"`
	Created_at  time.Time `json:"created_at"`
	Modified_at time.Time `json:"modified_at"`
}
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Change the function name to start with an uppercase letter
func GetMuxRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", processJSON).Methods("POST")
	r.HandleFunc("/processRequest", processRequestHandler).Methods("POST")
	r.HandleFunc("/processGRPCRequest", processGRPCRequest).Methods("POST")
	return r
}

func processJSON(w http.ResponseWriter, r *http.Request) {
	read, _ := io.ReadAll(r.Body)
	stringbody := string(read)
	res := JsonResponse{
		Error:   false,
		Message: "Reached Broker",
		Body:    stringbody,
	}
	bytes, _ := json.MarshalIndent(res, "", "\t")
	w.Write(bytes)
}
func processRequestHandler(w http.ResponseWriter, r *http.Request) {
	var brokerReq BrokerRequest
	err := json.NewDecoder(r.Body).Decode(&brokerReq)
	if err != nil {
		fmt.Println("got error: ", err)
		jsonres := &JsonResponse{
			Error:   true,
			Message: fmt.Sprintf("Got error: %v", err),
			Body:    fmt.Sprintf("Got error: %v", err),
		}
		json.NewEncoder(w).Encode(jsonres)
		return
	}
	action := brokerReq.Action
	switch strings.Split(action, "_")[0] {
	case "auth":
		handleAuthRequest(w, strings.Split(action, "_")[1], brokerReq.Payload)
	case "log":
		// handleLogRequest(w, strings.Split(action, "_")[1], brokerReq.LogPayload)
		handleLogRPCRequest(w, strings.Split(action, "_")[1], brokerReq.LogPayload)
	case "mail":
		handleMailRequest(w, strings.Split(action, "_")[1], brokerReq.MailPayload)
	default:
		jsonres := &JsonResponse{
			Error:   true,
			Message: "Pass valid request",
			Body:    "Pass valid request",
		}
		json.NewEncoder(w).Encode(jsonres)
		return

	}

}
func handleAuthRequest(w http.ResponseWriter, action string, payload AuthPayload) {

	auth_payload, _ := json.Marshal(payload)
	fmt.Println("going to call auth with url:" + "http://dockompose-authentication-service-1:8081/" + action)
	resp, err := http.Post("http://dockompose-authentication-service-1:8081/"+action, "application/json", bytes.NewReader(auth_payload))
	if err != nil {
		fmt.Println("got error in handleAuthRequest:", err)
		jsonres := &JsonResponse{
			Error:   true,
			Message: fmt.Sprintf("Got error: %v", err),
			Body:    fmt.Sprintf("Got error: %v", err),
		}
		json.NewEncoder(w).Encode(jsonres)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK status code:", resp.Status)
		// Handle non-OK status codes appropriately
		// Example: Return an error response
		errorMessage := fmt.Sprintf("Received non-OK status code: %v", resp.Status)
		jsonres := &JsonResponse{
			Error:   true,
			Message: errorMessage,
			Body:    errorMessage,
		}
		json.NewEncoder(w).Encode(jsonres)
		return // Exit the function
	}
	defer resp.Body.Close()

	var responsePayload AuthPayload
	json.NewDecoder(resp.Body).Decode(&responsePayload)
	json.NewEncoder(w).Encode(responsePayload)
}
func handleLogRequest(w http.ResponseWriter, action string, payload LogPayload) {

	log_payload, _ := json.Marshal(payload)
	fmt.Println("going to call auth with url:" + "http://dockompose-logger-service-1/" + action)
	resp, err := http.Post("http://dockompose-logger-service-1/"+action, "application/json", bytes.NewReader(log_payload))
	if err != nil {
		fmt.Println("got error in handleLogRequest:", err)
		jsonres := &JsonResponse{
			Error:   true,
			Message: fmt.Sprintf("Got error: %v", err),
			Body:    fmt.Sprintf("Got error: %v", err),
		}
		json.NewEncoder(w).Encode(jsonres)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK status code:", resp.Status)
		// Handle non-OK status codes appropriately
		// Example: Return an error response
		errorMessage := fmt.Sprintf("Received non-OK status code: %v", resp.Status)
		jsonres := &JsonResponse{
			Error:   true,
			Message: errorMessage,
			Body:    errorMessage,
		}
		json.NewEncoder(w).Encode(jsonres)
		return // Exit the function
	}
	defer resp.Body.Close()

	// var responsePayload LogPayload
	// json.NewDecoder(resp.Body).Decode(&responsePayload)
	fmt.Println("got response payload: ", resp.Body)
	bytes, _ := io.ReadAll(resp.Body)
	w.Write(bytes)
	// json.NewEncoder(w).Encode(resp.Body)
}

func handleMailRequest(w http.ResponseWriter, action string, payload MailPayload) {

	log_payload, _ := json.Marshal(payload)
	fmt.Println("going to call auth with url:" + "http://dockompose-mail-service-1/" + action)
	resp, err := http.Post("http://dockompose-mail-service-1/"+action, "application/json", bytes.NewReader(log_payload))
	if err != nil {
		fmt.Println("got error in handleLogRequest:", err)
		jsonres := &JsonResponse{
			Error:   true,
			Message: fmt.Sprintf("Got error: %v", err),
			Body:    fmt.Sprintf("Got error: %v", err),
		}
		json.NewEncoder(w).Encode(jsonres)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK status code:", resp.Status)
		// Handle non-OK status codes appropriately
		// Example: Return an error response
		errorMessage := fmt.Sprintf("Received non-OK status code: %v", resp.Status)
		jsonres := &JsonResponse{
			Error:   true,
			Message: errorMessage,
			Body:    errorMessage,
		}
		json.NewEncoder(w).Encode(jsonres)
		return // Exit the function
	}
	defer resp.Body.Close()

	// var responsePayload LogPayload
	// json.NewDecoder(resp.Body).Decode(&responsePayload)
	fmt.Println("got response payload: ", resp.Body)
	bytes, _ := io.ReadAll(resp.Body)
	w.Write(bytes)
	// json.NewEncoder(w).Encode(resp.Body)
}

type RPCPayload struct {
	Name string
	Data string
}

func handleLogRPCRequest(w http.ResponseWriter, action string, payload LogPayload) error {
	client, err := rpc.Dial("tcp", "dockompose-logger-service-1:5001")
	if err != nil {
		fmt.Println("got error", err)
	}
	rpcPayload := RPCPayload{
		Name: payload.Name,
		Data: payload.Data,
	}
	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		fmt.Println("got error", err)
	}

	return nil

}
func processGRPCRequest(w http.ResponseWriter, r *http.Request) {
	var logPayload LogPayload

	err := json.NewDecoder(r.Body).Decode(&logPayload)
	if err != nil {
		w.Write([]byte("Pass valid request"))
	}
	conn, err := grpc.Dial("dockompose-logger-service-1:5050", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		w.Write([]byte("Not able to dial Grpc"))
	}
	defer conn.Close()
	c := grpcproto.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.WriteLog(ctx, &grpcproto.LogRequest{
		LogEntry: &grpcproto.Log{
			Name: logPayload.Name,
			Data: logPayload.Data,
		},
	})
	if err != nil {
		w.Write([]byte("Got error" + err.Error()))
	}
	fmt.Println("successfully proccessed log entry")
	w.Write([]byte("logged entry"))
}
