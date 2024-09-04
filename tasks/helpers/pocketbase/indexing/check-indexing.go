package indexing_helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Arinji2/vibeify-backend/api"
	pocketbase_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/pocketbase"
	"github.com/Arinji2/vibeify-backend/types"
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
		return
	}

	jsonFile, error := os.ReadFile("/tasks/helpers/pocketbase/indexing/indexable-playlists.json")

	if error != nil {
		fmt.Println(error)
		return
	}

	jsonData := []types.IndexablePlaylist{}
	var pool = make(chan struct{}, 2)

	error = json.Unmarshal(jsonFile, &jsonData)
	if error != nil {
		fmt.Println(error)
		return
	}

	isIndexing := IsIndexingSongs()
	fmt.Println(isIndexing)
	if isIndexing {
		fmt.Println("Indexing is already in progress")
		return
	}

	fmt.Println("INDEXING PLAYLISTS")

	for _, playlist := range jsonData {

		pool <- struct{}{}
		go func(playlist types.IndexablePlaylist) {
			defer func() { <-pool }()

			PerformPlaylistIndexing(playlist)

		}(playlist)

	}

}
