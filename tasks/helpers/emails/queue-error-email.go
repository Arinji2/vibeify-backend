package email_helpers

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed templates/queue-error.html
var queueErrorEmailTemplateString string

func SendQueueErrorEmail(errorMsg string, email string) {
	type emailDataType struct {
		ErrorMsg string
		Year     int
	}

	emailData := emailDataType{
		ErrorMsg: errorMsg,
		Year:     time.Now().Year(),
	}

	emailString := emailTemplateUtility(emailData, "Queue Error Email", queueErrorEmailTemplateString)
	fmt.Println(emailString)

}
