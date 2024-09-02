package pocketbase_helpers

import (
	"fmt"
	"sync"
	"time"

	"github.com/Arinji2/vibeify-backend/api"
	"github.com/Arinji2/vibeify-backend/types"
)

func RecordPlaylistForDeletion(user *types.PocketbaseUser, playlists []types.SpotifyPlaylist) (errorText string) {

	adminToken, err := GetPocketbaseAdminToken()
	errorText = "Server Error"
	if err != "" {
		fmt.Println(err)
		return
	}
	client := api.NewApiClient()

	var dateToBeDeleted string
	if user.Record.Premium {
		deletionTime := time.Now().Add(time.Hour * 50)
		dateToBeDeleted = deletionTime.Format("2006-01-02 15:04:05.000Z")
	} else {
		deletionTime := time.Now().Add(time.Hour * 24)
		dateToBeDeleted = deletionTime.Format("2006-01-02 15:04:05.000Z")
	}
	var wg sync.WaitGroup

	for _, playlist := range playlists {
		wg.Add(1)
		go func(playlist types.SpotifyPlaylist) {
			defer wg.Done()
			client.SendRequestWithBody("POST", "/api/collections/convertDeletion/records", map[string]string{
				"toBeDeleted": dateToBeDeleted,
				"playlistID":  playlist.ID,
			}, map[string]string{
				"Authorization": adminToken,
			})

		}(playlist)
	}

	wg.Wait()
	errorText = ""
	return
}
