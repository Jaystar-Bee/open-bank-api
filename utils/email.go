package utils

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/jordan-wright/email"
	"github.com/mailersend/mailersend-go"
)

func SendEmail(to, subject, body string, to_name string, tags []string, status chan bool) {
	e := email.NewEmail()
	e.From = fmt.Sprintf("Open Bank <%s>", os.Getenv("EMAIL_ACCOUNT"))
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(body)

	var err error
	if os.Getenv("ENV") == "development" {
		err = e.Send("localhost:1025", nil)
	} else {
		ms := mailersend.NewMailersend(os.Getenv("MAILERSEND_API_KEY"))
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		from := mailersend.From{
			Name:  "Open Bank",
			Email: os.Getenv("EMAIL_ACCOUNT"),
		}
		recipients := []mailersend.Recipient{
			{
				Name:  to_name,
				Email: to,
			},
		}

		message := ms.Email.NewMessage()

		message.SetFrom(from)
		message.SetRecipients(recipients)
		message.SetSubject(subject)
		message.SetHTML(body)
		message.SetTags(tags)
		message.SetInReplyTo("client-id")

		_, err = ms.Email.Send(ctx, message)
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
