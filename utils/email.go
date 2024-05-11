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

	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", os.Getenv("EMAIL_ACCOUNT"), os.Getenv("EMAIL_PASSWORD"), "smtp.gmail.com"))

	// err := e.Send("localhost:1025", smtp.PlainAuth("", os.Getenv("EMAIL_ACCOUNT"), os.Getenv("EMAIL_PASSWORD"), "localhost"))
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
