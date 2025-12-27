package main

import (
	"bookings/internals/models"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMessage(msg)
		}
	}()
}

func sendMessage(m models.MailData) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mail_username := os.Getenv("MAIL_USERNAME")
	mail_password := os.Getenv("MAIL_PASSWORD")
	ml := gomail.NewMessage()
	ml.SetHeader("From", m.From)
	ml.SetHeader("To", m.To)
	ml.SetHeader("Subject", m.Subject)
	if m.Template == "" {
		ml.SetBody("text/html", m.Content)
	} else {
		data, err := os.ReadFile(fmt.Sprintf("./templates/email-template/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}
		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		ml.SetBody("text/html", msgToSend)
	}
	d := gomail.NewDialer("smtp.gmail.com", 587, mail_username, mail_password)
	if err := d.DialAndSend(ml); err != nil {
		panic(err)
	}
}
