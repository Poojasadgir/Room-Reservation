package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/Poojasadgir/room-reservation/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

// listenForMail listens for incoming messages on the app.MailChannel and sends them using the sendMessage function.
func listenForMail() {
	go func() {
		for {
			message := <-app.MailChannel
			sendMessage(message)
		}
	}()
}

// sendMessage sends an email using the provided MailData.
// If a template is provided, it will be used to format the email content.
// Otherwise, the content will be used as is.
func sendMessage(m models.MailData) {
	const mailTimeout = 10 * time.Second
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = mailTimeout
	server.SendTimeout = mailTimeout

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}
		mailTemplate := string(data)
		messageToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, messageToSend)
	}
	err = email.Send(client)
	if err != nil {
		errorLog.Println(err)
	} else {
		infoLog.Println("Email sent!")
	}
}
