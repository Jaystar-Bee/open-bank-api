package emails

import (
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

func SendOTPToMail(receiver, name, otp string) error {
	e := email.NewEmail()
	e.From = "<" + os.Getenv("EMAIL_ACCOUNT") + ">"
	e.To = []string{"jbayilara@gmail.com"}
	e.Subject = "Awesome Subject"
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", os.Getenv("EMAIL_ACCOUNT"), os.Getenv("EMAIL_PASSWORD"), "smtp.gmail.com"))
	if err != nil {
		panic(err)
	}
	return err

}
