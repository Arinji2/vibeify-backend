package spotify_helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Arinji2/vibeify-backend/api"
	"github.com/Arinji2/vibeify-backend/types"
)

var (
	retries int = 0
)

func GetSpotifyPlaylist(url string, user *types.PocketbaseUser) (err error) {
	token := GetSpotifyToken()

	playlistID := strings.Split(strings.Split(url, "/")[4], "?")[0]

	client := api.NewApiClient("https://api.spotify.com/v1")
	res, status, err := client.SendRequestWithQuery("GET", fmt.Sprintf("/playlists/%s", playlistID), map[string]string{

		"fields": "id,name,description,owner(display_name,id),tracks.items(added_at,track(id,name,artists(id,name),album(id,name),external_urls.spotify,uri)),tracks.total,tracks.offset,tracks.limit",
	}, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	})

	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if status == 400 {
		err = errors.New("playlist not found")
		return
	}

	if status == 401 {
		BustSpotifyTokenCache()
		retries++
		if retries < 3 {
			fmt.Println("Retrying Spotify Authentication")
			GetSpotifyPlaylist(url, user)
		} else {
			err = errors.New("spotify token expired")
			return
		}
	}

	var Playlist types.SpotifyPlaylist
	jsonData, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("Error marshalling response: %v", err)
	}

	err = json.Unmarshal(jsonData, &Playlist)
	if err != nil {
		log.Fatalf("Error unmarshalling into struct: %v", err)
	}

	if Playlist.Tracks.Total > 200 && !user.Record.Premium {
		err = errors.New("playlist is too large for free users")
	}

	if Playlist.Tracks.Total > 400 && user.Record.Premium {
		err = errors.New("playlist is too large for premium users")
	}

	fmt.Println(Playlist.Tracks.Total, Playlist.Tracks.Limit, Playlist.Tracks.Offset)
	return

}
