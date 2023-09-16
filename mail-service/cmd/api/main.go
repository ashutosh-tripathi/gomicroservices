package main

import (
	"fmt"
	"mail-service/cmd/api/router"
	"net/http"
)

const webPort = 80

func main() {

	r := router.GetRouter("dockompose-mailhog-1", 1025)

	defer http.ListenAndServe(fmt.Sprintf(":%d", webPort), r)

	// mailSer := &mailer.MessageServer{}

	// msg := &mailer.Message{
	// 	From:    "rbj.ashu@gmail.com",
	// 	To:      "mitra.rituparna04@gmail.com",
	// 	Subject: "Hello",
	// 	Data:    "I LUV U",
	// }
	// err := mailSer.SendMessage(msg)
	// if err != nil {
	// 	fmt.Println("Got error: ", err)
	// }
	// fmt.Println("successfully sent message")

}
