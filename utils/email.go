package utils

import (
	"log"
	"os"

	"github.com/wneessen/go-mail"
)

func SendEmail(to, subject, body string, to_name string, tags []string, status chan bool) {

	m := mail.NewMsg()
	if err := m.FromFormat("Open Bank", os.Getenv("GOOGLE_MAIL_EMAIL")); err != nil {
		log.Printf("failed to set From address: %s", err)
	}
	if err := m.To(to); err != nil {
		log.Printf("failed to set To address: %s", err)
	}
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextHTML, body)

	// I'm uusing gmail client
	c, err := mail.NewClient("smtp.gmail.com",
		mail.WithPort(587), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(os.Getenv("GOOGLE_MAIL_EMAIL")), mail.WithPassword(os.Getenv("GOOGLE_MAIL_APP_PASSWORD")), mail.WithTLSPolicy(mail.TLSMandatory))

	if err != nil {
		log.Printf("failed to create mail client: %s", err)
		status <- false
		return
	}

	// Send the mail
	if err := c.DialAndSend(m); err != nil {
		log.Printf("failed to send mail: %s", err)
		status <- false
		return
	}
	status <- true
}
