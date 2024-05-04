package utils

import (
	"os"
	"time"

	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string, status chan bool) {
	// set up email
	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("EMAIL_ACCOUNT"))
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)
	// send email
	data := gomail.NewDialer("smtp.example.com", 587, os.Getenv("EMAIL_ACCOUNT"), os.Getenv("EMAIL_PASSWORD"))

	var err error

	for i := 0; i < 3; i++ {
		err = data.DialAndSend(message)
		if err == nil {
			status <- true
			return
		}
		time.Sleep(time.Second * 4)
	}

	status <- false

}
