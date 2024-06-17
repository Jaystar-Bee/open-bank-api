package emails

import (
	"errors"
	"fmt"
	"os"

	"github.com/Jaystar-Bee/open-bank-api/utils"
)

func SendWelcomeEmail(to_email string, to_name string) error {
	emailData := map[string]string{
		"Name":      to_name,
		"Login_url": fmt.Sprintf("%s/auth/login", os.Getenv("WEB_URL")),
	}
	body, err := utils.ParseTemplate("emails/templates/welcome.html", emailData)
	if err != nil {
		return err
	}
	emailCompleted := make(chan bool)
	go utils.SendEmail(to_email, "Welcome to Open Bank", body, to_name, []string{"onboarding", "welcome"}, emailCompleted)

	if <-emailCompleted {
		return nil
	} else {
		return errors.New("unable to send email")
	}

}
