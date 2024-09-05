package pocketbase_helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/Arinji2/vibeify-backend/api"
	"github.com/Arinji2/vibeify-backend/types"
	user_errors "github.com/Arinji2/vibeify-backend/user-errors"
)

func GetInternalGenre(tracks []types.SpotifyPlaylistItem, genres []string, genreArrays types.GenreArrays) (updatedTracks []types.SpotifyPlaylistItem, err error) {
	wg := sync.WaitGroup{}
	AdminToken, err := GetPocketbaseAdminToken()
	if err != nil {
		return nil, user_errors.NewUserError("", err)
	}
	client := api.NewApiClient()

	genreString := strings.Builder{}
	for i, genre := range genres {
		genreString.WriteString(fmt.Sprintf(`genres ~ "%s"`, genre))
		if i != len(genres)-1 {
			genreString.WriteString(" || ")
		}

	}

	errorChannel := make(chan error, len(tracks))

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
				errorChannel <- user_errors.NewUserError("", err)
				return
			}

			totalItems, ok := res["totalItems"].(float64)
			if !ok {
				errorChannel <- user_errors.NewUserError("", errors.New("totalItems is not a float64"))
				return
			}
			if totalItems == 0 {
				updatedTracks = append(updatedTracks, track)
				return
			}

			items, ok := res["items"].([]interface{})
			if !ok {
				errorChannel <- user_errors.NewUserError("", errors.New("items is not a []interface{}"))
				return
			}

			marshalledItems, err := json.Marshal(items)
			if err != nil {
				errorChannel <- user_errors.NewUserError("", err)
				return
			}

			records := []types.PocketbaseSongRecord{}
			err = json.Unmarshal(marshalledItems, &records)
			if err != nil {

				errorChannel <- user_errors.NewUserError("", err)
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
						URI: track.Track.URI,
					})

				}

			}

			if !hasMatched {
				updatedTracks = append(updatedTracks, track)

			}

			return

		}(track, genres)

		wg.Wait()
		close(errorChannel)
	}

	var finalError error
	for err := range errorChannel {
		if err != nil {
			if finalError == nil {
				finalError = err
			} else {
				finalError = fmt.Errorf("%v; %v", finalError, err)
			}
		}
	}

	if finalError != nil {
		return nil, user_errors.NewUserError("", finalError)
	}

	return updatedTracks, nil

}
