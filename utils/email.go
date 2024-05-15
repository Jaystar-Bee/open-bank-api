package utils

import (
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

func SendEmail(to, subject, body string, status chan bool) {
	e := email.NewEmail()
	e.From = fmt.Sprintf("Open Bank <%s>", os.Getenv("EMAIL_ACCOUNT"))
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(body)

	var err error
	if os.Getenv("ENV") == "development" {
		err = e.Send("localhost:1025", nil)
	} else {
		err = e.Send(fmt.Sprint(os.Getenv("MAIL_HOST"), ":587"), smtp.PlainAuth("", os.Getenv("MAIL_USERNAME"), os.Getenv("SMT_TOKEN"), os.Getenv("MAIL_HOST")))
	}

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		if smtpErr, ok := err.(net.Error); ok && smtpErr.Timeout() {
			log.Println("Temporary SMTP error, retrying later...")
		} else {
			log.Println("Permanent SMTP error, cannot retry.")
		}
		status <- false
		return // Return the error instead of exiting
	}
	status <- true
	fmt.Println("Email sent successfully!")
}
