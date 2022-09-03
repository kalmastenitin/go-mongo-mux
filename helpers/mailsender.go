package helpers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func GetCreds() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading Env File")
	}
	return os.Getenv("EMAILKEY")
}

func CheckEmail(usermail, username string) bool {
	log.Println(usermail, username)
	from := mail.NewEmail("goapptest", "nk@diycam.co.in")
	subject := "User Checking User Email"
	to := mail.NewEmail(username, usermail)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>You can now use your account</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(GetCreds())
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return false
	} else {
		log.Println(response.StatusCode)
		log.Println(response.Body)
		log.Println(response.Headers)
		if response.StatusCode == 202 {
			return true
		}
		return false
	}
}
