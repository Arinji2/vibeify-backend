package compare

import (
	"fmt"
	"strings"

	"github.com/Arinji2/vibeify-backend/api"

	"github.com/Arinji2/vibeify-backend/types"
)

func CheckExistingCompares(user types.PocketbaseUser) bool {
	client := api.NewApiClient()
	res, _, err := client.SendRequestWithQuery("GET", "/api/collections/compareList/records", map[string]string{
		"page":    "1",
		"perPage": "500",
		"filter":  fmt.Sprintf(`user="%s" && results=null`, user.Record.ID),
	}, map[string]string{
		"Authorization": user.Token,
	})

	if err != nil {
		fmt.Println(err)
		return false
	}

	totalItems, ok := res["totalItems"].(float64)
	if !ok {
		fmt.Println("Error: Total Results not found")
		return false
	}

	if totalItems > 0 {
		return true
	}

	return false

}

func CheckIfCompareExists(user types.PocketbaseUser, taskData types.CompareTaskType) bool {
	playlist1ID := strings.Split(strings.Split(taskData.Playlist1, "/")[4], "?")[0]
	playlist2ID := strings.Split(strings.Split(taskData.Playlist2, "/")[4], "?")[0]

	client := api.NewApiClient()
	res, _, err := client.SendRequestWithQuery("GET", "/api/collections/compareList/records", map[string]string{
		"page":    "1",
		"perPage": "1",
		"filter":  fmt.Sprintf(`user="%s" && playlist1="%s" && playlist2="%s"`, user.Record.ID, playlist1ID, playlist2ID),
	}, map[string]string{
		"Authorization": user.Token,
	})

	if err != nil {
		fmt.Println("Error: Compare not found")
		return false
	}

	totalItems, ok := res["totalItems"].(float64)
	if !ok {
		fmt.Println("Error: Total Results not found")
		return false
	}

	if totalItems > 0 {
		return true
	}

	return false
}
