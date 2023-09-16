package mailer

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

type MessageServer struct {
	SMTPHOST     string
	SMTPPORT     int
	SMTPUSER     string
	SMTPPASSWORD string
}
type Message struct {
	From    string
	To      string
	Data    string
	Subject string
}

func (m *MessageServer) SendMessage(msg *Message) error {
	mes := gomail.NewMessage()
	mes.SetHeader("From", msg.From)
	mes.SetHeader("To", msg.To)
	mes.SetHeader("Subject", msg.Subject)
	mes.SetBody("text/plain", msg.Data)

	d := gomail.NewDialer(m.SMTPHOST, m.SMTPPORT, m.SMTPUSER, m.SMTPPASSWORD)

	// Send the email
	if err := d.DialAndSend(mes); err != nil {
		log.Println(err)
		return err
	} else {
		fmt.Println("Email sent successfully!")
	}
	return nil
}
