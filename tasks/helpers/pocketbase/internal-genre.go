package pocketbase_helpers

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/Arinji2/vibeify-backend/api"
	"github.com/Arinji2/vibeify-backend/types"
)

func GetInternalGenre(tracks []types.SpotifyPlaylistItem, genres []string, genreArrays types.GenreArrays) (errorString string, updatedTracks []types.SpotifyPlaylistItem) {
	wg := sync.WaitGroup{}
	AdminToken, errorString := GetPocketbaseAdminToken()
	if errorString != "" {
		return
	}
	client := api.NewApiClient()

	genreString := strings.Builder{}
	for i, genre := range genres {
		genreString.WriteString(fmt.Sprintf(`genres ~ "%s"`, genre))
		if i != len(genres)-1 {
			genreString.WriteString(" || ")
		}

	}

	errorString = "Server Error"

	for _, track := range tracks {
		wg.Add(1)
		go func(track types.SpotifyPlaylistItem, genres []string) {
			defer wg.Done()
			res, _, err := client.SendRequestWithQuery("GET", "/api/collections/songs/records", map[string]string{
				"page":    "1",
				"perPage": "1",
				"filter":  fmt.Sprintf(`spotifyID = "%s" && (%s)`, track.Track.ID, genreString.String()),
			}, map[string]string{
				"Authorization": AdminToken,
			})
			if err != nil {
				fmt.Println("Error in getting song record", err)
				return
			}

			totalItems, ok := res["totalItems"].(float64)
			if !ok {
				fmt.Println("Error converting totalItems to int:", err)
				return
			}
			if totalItems == 0 {
				updatedTracks = append(updatedTracks, track)
				return
			}

			items, ok := res["items"].([]interface{})
			if !ok {
				fmt.Println("Error converting items to []interface{}:", err)
				return
			}

			marshalledItems, err := json.Marshal(items)
			if err != nil {
				fmt.Println("Error marshalling items:", err)
				return
			}

			records := []types.PocketbaseSongRecord{}
			err = json.Unmarshal(marshalledItems, &records)
			if err != nil {
				fmt.Println("Error unmarshalling records:", err)
				return
			}

			record := records[0]
			hasMatched := false

			for i, genre := range record.Genres {
				record.Genres[i] = strings.ToLower(genre)

			}

			for _, genre := range genres {

				genreMatch := slices.Contains(record.Genres, genre)
				if genreMatch {
					hasMatched = true

					genreArrays[genre] = append(genreArrays[genre], types.GenreArray{
						URI:  track.Track.URI,
						Name: track.Track.Name,
					})

				}

			}

			if !hasMatched {
				updatedTracks = append(updatedTracks, track)

			}

		}(track, genres)

		wg.Wait()

	}

	errorString = ""

	return
}
