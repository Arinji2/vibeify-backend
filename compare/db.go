package compare

import (
	"fmt"
	"strings"

	"github.com/Arinji2/vibeify-backend/api"
	custom_log "github.com/Arinji2/vibeify-backend/logger"
	pocketbase_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/pocketbase"
	"github.com/Arinji2/vibeify-backend/types"
)

func AddToDb(user types.PocketbaseUser, taskData types.CompareTaskType) bool {
	client := api.NewApiClient()
	adminToken := pocketbase_helpers.GetPocketbaseAdminToken()
	playlist1ID := strings.Split(strings.Split(taskData.Playlist1, "/")[4], "?")[0]
	playlist2ID := strings.Split(strings.Split(taskData.Playlist2, "/")[4], "?")[0]

	_, _, err := client.SendRequestWithBody("POST", "/api/collections/compareList/records", map[string]string{
		"playlist1": playlist1ID,
		"playlist2": playlist2ID,
		"user":      user.Record.ID,
	}, map[string]string{
		"Authorization": adminToken,
	})

	if err != nil {
		fmt.Println("Error: Adding to DB")
		return false
	}

	custom_log.Logger.Debug("Added New Compare to DB")
	return true
}
