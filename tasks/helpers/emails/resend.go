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

func EmailClient(to string, subject string) *EmailClientType {
	godotenv.Load()
	apiKey := os.Getenv("EMAIL_KEY")
	emailClient := resend.NewClient(apiKey)

	emailDetails := &resend.SendEmailRequest{
		From:    "no-reply@mail.arinji.com",
		To:      []string{to},
		Subject: subject,
	}

	returnData := EmailClientType{
		Client:     emailClient,
		SendParams: emailDetails,
	}

	return &returnData

}

func (e *EmailClientType) SendEmail(subject string, html string) {

	e.SendParams.Subject = subject
	e.SendParams.Text = html
	_, err := e.Client.Emails.Send(e.SendParams)
	if err != nil {
		panic(err)
	}
}
