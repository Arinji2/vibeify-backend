package tasks

import (
	"fmt"

	"github.com/Arinji2/vibeify-backend/tasks/helpers"
	pocketbase_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/pocketbase"
	"github.com/Arinji2/vibeify-backend/types"
)

func PerformTask(task types.AddTaskType) {
	user, err := pocketbase_helpers.ValidateUser(task.UserToken)
	if err != nil {
		helpers.HandleError(err, "")
	}

	used, total, err := pocketbase_helpers.CheckLimit(user)
	fmt.Println(used, total)
	if err != nil {
		helpers.HandleError(err, user.Record.Email)
	}

	//email_helpers.SendQueueAdditionEmail(user.Record.Premium, user.Record.Email)
}
