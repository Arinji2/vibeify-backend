package indexing_helpers

import (
	"fmt"
	"sync"

	"github.com/Arinji2/vibeify-backend/api"
	pocketbase_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/pocketbase"
	"github.com/Arinji2/vibeify-backend/types"
)

func QueueSongIndexing(tracks []types.SpotifyPlaylistItem, priorityIndex string) {
	adminToken, err := pocketbase_helpers.GetPocketbaseAdminToken()
	if err != "" {
		fmt.Println(err)
		return
	}

	client := api.NewApiClient()
	var wg sync.WaitGroup
	pool := make(chan struct{}, 10)
	defer close(pool)

	for _, track := range tracks {
		wg.Add(1)
		pool <- struct{}{}
		go func(track types.SpotifyPlaylistItem) {
			defer wg.Done()
			defer func() { <-pool }()

			if err := sendSongToIndex(client, adminToken, track.Track.ID, priorityIndex); err != nil {
				fmt.Println(err)
			}
		}(track)
	}

	wg.Wait()
}
