package indexing_helpers

import (
	"fmt"
	"sync"

	"github.com/Arinji2/vibeify-backend/api"
	pocketbase_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/pocketbase"
	"github.com/Arinji2/vibeify-backend/types"
)

func PerformSongIndexing() {
	client := api.NewApiClient()
	adminToken, err := pocketbase_helpers.GetPocketbaseAdminToken()
	if err != "" {
		fmt.Println(err)
		return
	}

	songsToIndex, error := getSongsToIndex(client, adminToken)
	if error != nil {
		fmt.Println(err)
		return
	}

	spotifyTracks, error := fetchSpotifyTracks(songsToIndex)
	if error != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	pool := make(chan struct{}, 2)

	for _, song := range spotifyTracks {
		wg.Add(1)
		pool <- struct{}{}
		go func(song types.SpotifyTrack) {
			defer wg.Done()
			defer func() { <-pool }()

			if err := indexSong(client, adminToken, song); err != nil {
				fmt.Println(err)
			}
		}(song)
	}

	wg.Wait()
}
