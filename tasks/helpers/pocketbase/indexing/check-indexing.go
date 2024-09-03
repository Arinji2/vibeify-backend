package indexing_helpers

import (
	"fmt"
	"sync"

	"github.com/Arinji2/vibeify-backend/api"
	pocketbase_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/pocketbase"
)

var (
	inProgress sync.Mutex
)

func CheckIndexing() {
	inProgress.Lock()
	defer inProgress.Unlock()

	client := api.NewApiClient()
	adminToken, err := pocketbase_helpers.GetPocketbaseAdminToken()
	if err != "" {
		fmt.Println(err)
		return
	}

	res, _, error := client.SendRequestWithQuery("GET", "/api/collections/songsToIndex/records", map[string]string{
		"page":    "1",
		"perPage": "1",
		"fields":  "id"}, map[string]string{
		"Authorization": adminToken,
	})

	if error != nil {
		fmt.Println(err)
		return
	}

	totalItems, ok := res["totalItems"].(float64)
	if !ok {
		fmt.Println("Error getting total items")
	}

	if totalItems > 0 {
		PerformSongIndexing()
	}

	//TODO Read a json file to get playlists to keep indexed
}
