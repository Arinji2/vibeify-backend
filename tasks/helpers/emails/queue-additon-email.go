package email_helpers

import (
	_ "embed"
	"fmt"
)

//go:embed templates/queue-added.html
var queueAdditionEmailTemplateString string

func SendQueueAdditionEmail(isPremium bool, email string) {
	type emailDataType struct {
		IsPremium bool
	}

	emailData := emailDataType{
		IsPremium: isPremium,
	}

	emailString := emailTemplateUtility(emailData, "Queue Addition Email", queueAdditionEmailTemplateString)
	fmt.Println(emailString)

}
