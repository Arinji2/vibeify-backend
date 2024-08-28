package email_helpers

import (
	_ "embed"
	"fmt"
)

//go:embed templates/queue-error.html
var queueErrorEmailTemplateString string

func SendQueueErrorEmail(errorMsg string, email string) {
	type emailDataType struct {
		ErrorMsg string
	}

	emailData := emailDataType{
		ErrorMsg: errorMsg,
	}

	emailString := emailTemplateUtility(emailData, "Queue Error Email", queueErrorEmailTemplateString)
	fmt.Println(emailString)

}
