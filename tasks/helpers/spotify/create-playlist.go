package spotify_helpers

import (
	"fmt"
	"sync"

	"github.com/Arinji2/vibeify-backend/types"
)

func CreatePlaylists(playlistName string, genreArrays types.GenreArrays) (errorString string, createdPlaylists []types.SpotifyPlaylist) {

	errorString = "Server Error"
	accessToken := GetSpotifyToken()

	spotifyHeaders := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
	}
	var wg sync.WaitGroup

	for key, genre := range genreArrays {
		wg.Add(1)
		go func(key string, genre []types.GenreArray) {

			defer wg.Done()
			if len(genre) == 0 {
				return
			}
			err, createdPlaylist := createGenrePlaylist(key, playlistName, spotifyHeaders)
			if err != "" {

				fmt.Println(err)

				return
			}

			err = initPlaylist(createdPlaylist, genre, spotifyHeaders)
			if err != "" {

				fmt.Println(err)

				return
			}

			createdPlaylists = append(createdPlaylists, createdPlaylist)
		}(key, genre)
	}

	wg.Wait()

	if len(createdPlaylists) == 0 {
		errorString = "No playlists created"
		return
	}

	errorString = ""

	return
}
