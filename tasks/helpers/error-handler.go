package helpers

import (
	"fmt"

	email_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/emails"
)

func HandleError(err string, emailTo string) {
	fmt.Println("Handling Error")
	defer panic(err)
	if emailTo == "" {
		fmt.Println("No Email to Send To")
		return
	}
	fmt.Println("Sending Email to ", emailTo)

	email_helpers.SendQueueErrorEmail(err, emailTo)

}
