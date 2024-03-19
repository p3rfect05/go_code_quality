package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/p3rfect05/go_proj/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenToMail() {
	go func() {
		for {
			msg := <-appConfig.MailChannel
			sendMessage(msg)
		}
	}()
}

func sendMessage(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		appConfig.ErrorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := os.ReadFile(fmt.Sprintf("./email_templates/%s", m.Template))
		if err != nil {
			appConfig.ErrorLog.Println(err)
		}
		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}
	err = email.Send(client)
	if err != nil {
		appConfig.ErrorLog.Println(err)
	} else {
		appConfig.InfoLog.Println("Email sent successfully")
	}

}
