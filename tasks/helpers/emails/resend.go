package email_helpers

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/resend/resend-go/v2"
)

type EmailClientType struct {
	Client     *resend.Client
	SendParams *resend.SendEmailRequest
}

func NewEmailClient(to string) *EmailClientType {
	godotenv.Load()
	apiKey := os.Getenv("EMAIL_KEY")
	emailClient := resend.NewClient(apiKey)

	emailDetails := &resend.SendEmailRequest{
		From: "no-reply@mail.arinji.com",
		To:   []string{to},
	}

	returnData := EmailClientType{
		Client:     emailClient,
		SendParams: emailDetails,
	}

	return &returnData

}

func (e *EmailClientType) SendEmail(subject string, html string) {

	e.SendParams.Subject = subject
	e.SendParams.Html = html
	_, err := e.Client.Emails.Send(e.SendParams)
	if err != nil {
		panic(err)
	}
}
