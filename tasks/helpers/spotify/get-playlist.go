package spotify_helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Arinji2/vibeify-backend/api"
	"github.com/Arinji2/vibeify-backend/cache"
	"github.com/Arinji2/vibeify-backend/types"
)

var (
	retries       int = 0
	playlistCache     = cache.NewCache(500, 10*time.Minute)
)

func GetSpotifyPlaylist(url string, user *types.PocketbaseUser, indexingFlag ...bool) (tracks []types.SpotifyPlaylistItem, playlistName string, errorText string) {
	tracks = nil
	fmt.Println(strings.Split(url, "/"))
	playlistID := strings.Split(strings.Split(url, "/")[4], "?")[0]

	if cachedData, found := playlistCache.Get(playlistID); found {
		playlist := cachedData.(types.SpotifyPlaylist)
		tracks = playlist.Tracks.Items
		playlistName = playlist.Name
		return
	}

	token := GetSpotifyToken()

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
		errorText = "Playlist not found"
		return
	}

	if status == 401 {
		BustSpotifyTokenCache()
		retries++
		if retries < 3 {
			fmt.Println("Retrying Spotify Authentication")
			return GetSpotifyPlaylist(url, user)
		} else {
			fmt.Println("Spotify Token Expired")
			errorText = "Server Error"
			return
		}
	}

	var Playlist types.SpotifyPlaylist
	jsonData, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("Error marshalling response: %v", err)
		errorText = "Server Error"
		return
	}

	err = json.Unmarshal(jsonData, &Playlist)
	if err != nil {
		fmt.Printf("Error unmarshalling into struct: %v", err)
		errorText = "Server Error"
		return
	}

	isIndexing := false
	if len(indexingFlag) > 0 {
		isIndexing = indexingFlag[0]
	}
	if !isIndexing {
		if Playlist.Tracks.Total > 200 && !user.Record.Premium {
			errorText = "Playlist is too large. Maximum size is 200 tracks for free users"
			return
		}

		if Playlist.Tracks.Total > 400 && user.Record.Premium {
			errorText = "Playlist is too large. Maximum size is 400 tracks for premium users"
			return
		}
	}

	for {
		if (Playlist.Tracks.Offset + Playlist.Tracks.Limit) >= Playlist.Tracks.Total {
			break
		}
		Playlist.Tracks.Offset += Playlist.Tracks.Limit
		res, _, err = client.SendRequestWithQuery("GET", fmt.Sprintf("/playlists/%s/tracks", playlistID), map[string]string{
			"fields": "items(added_at,track(id,name,artists(id,name),album(id,name),external_urls.spotify,uri)),total,offset,limit",
			"limit":  "100",
			"offset": fmt.Sprintf("%d", Playlist.Tracks.Offset),
		}, map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", token),
		})

		Items := []types.SpotifyPlaylistItem{}
		if err != nil {
			fmt.Printf("Error sending request: %v", err)
			errorText = "Server Error"
			return
		}

		jsonData, err := json.Marshal(res["items"])
		if err != nil {
			fmt.Printf("Error marshalling response: %v", err)
			errorText = "Server Error"
			return
		}

		err = json.Unmarshal(jsonData, &Items)
		if err != nil {
			fmt.Printf("Error unmarshalling into struct: %v", err)
			errorText = "Server Error"
			return
		}

		Playlist.Tracks.Items = append(Playlist.Tracks.Items, Items...)
	}

	for _, track := range Playlist.Tracks.Items {
		if !track.Track.IsLocal {
			tracks = append(tracks, track)
		}
	}

	playlistCache.Set(playlistID, Playlist, time.Hour)

	playlistName = Playlist.Name

	return
}
