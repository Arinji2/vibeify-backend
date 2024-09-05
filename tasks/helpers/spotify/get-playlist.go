package spotify_helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Arinji2/vibeify-backend/api"
	"github.com/Arinji2/vibeify-backend/cache"
	"github.com/Arinji2/vibeify-backend/constants"
	"github.com/Arinji2/vibeify-backend/types"
	user_errors "github.com/Arinji2/vibeify-backend/user-errors"
)

var (
	retries       int = 0
	playlistCache     = cache.NewCache(500, 10*time.Minute)
)

func GetSpotifyPlaylist(url string, user *types.PocketbaseUser, indexingFlag ...bool) (tracks []types.SpotifyPlaylistItem, playlistName string, err error) {
	tracks = nil

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
		return nil, "", user_errors.NewUserError("", fmt.Errorf("error sending request: %v", err))
	}

	if status == 400 {

		return nil, "", user_errors.NewUserError("Playlist not found or Playlist is private", err)
	}

	if status == 401 {
		BustSpotifyTokenCache()
		retries++
		if retries < 3 {
			fmt.Println("Retrying Spotify Authentication")
			return GetSpotifyPlaylist(url, user)
		} else {

			return nil, "", user_errors.NewUserError("", errors.New("spotify token expired"))
		}
	}

	var Playlist types.SpotifyPlaylist
	jsonData, err := json.Marshal(res)
	if err != nil {

		return nil, "", user_errors.NewUserError("", (fmt.Errorf("error marshalling response: %v", err)))
	}

	err = json.Unmarshal(jsonData, &Playlist)
	if err != nil {
		return nil, "", user_errors.NewUserError("", fmt.Errorf("error unmarshalling into struct: %v", err))
	}

	isIndexing := false
	if len(indexingFlag) > 0 {
		isIndexing = indexingFlag[0]
	}
	if !isIndexing {
		if Playlist.Tracks.Total > constants.MAX_FREE_PLAYLIST_SIZE && !user.Record.Premium {
			return nil, "", user_errors.NewUserError(fmt.Sprintf("Playlist is too large. Maximum size is %d tracks for free users", constants.MAX_FREE_PLAYLIST_SIZE), errors.New("playlist-exceeds-free-limit"))
		}

		if Playlist.Tracks.Total > constants.MAX_PAID_PLAYLIST_SIZE && user.Record.Premium {
			return nil, "", user_errors.NewUserError(fmt.Sprintf("Playlist is too large. Maximum size is %d tracks for premium users", constants.MAX_PAID_PLAYLIST_SIZE), errors.New("playlist-exceeds-paid-limit"))
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
			return nil, "", user_errors.NewUserError("", fmt.Errorf("error sending request: %v", err))
		}

		jsonData, err := json.Marshal(res["items"])
		if err != nil {
			return nil, "", user_errors.NewUserError("", fmt.Errorf("error marshalling response: %v", err))
		}

		err = json.Unmarshal(jsonData, &Items)
		if err != nil {
			return nil, "", user_errors.NewUserError("", fmt.Errorf("error unmarshalling into struct: %v", err))
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

	return tracks, playlistName, nil
}
