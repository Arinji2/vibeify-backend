package tasks

import (
	"fmt"
	"strings"

	"github.com/Arinji2/vibeify-backend/tasks/helpers"
	ai_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/ai"
	email_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/emails"
	pocketbase_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/pocketbase"
	spotify_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/spotify"
	"github.com/Arinji2/vibeify-backend/types"
)

func PerformTask(task types.AddTaskType) {
	user, err := pocketbase_helpers.ValidateUser(task.UserToken)
	if err != "" {
		helpers.HandleError(err, "")
	}

	used, usesID, err := pocketbase_helpers.CheckLimit(user)

	if err != "" {
		helpers.HandleError(err, user.Record.Email)
	}

	email_helpers.SendQueueAdditionEmail(user.Record.Premium, user.Record.Email)
	spotify_helpers.GetSpotifyToken()

	tracks, playlistName, err := spotify_helpers.GetSpotifyPlaylist(task.SpotifyURL, user)
	if err != "" {
		helpers.HandleError(err, user.Record.Email)
	}

	genreArrays := types.GenreArrays{}

	for i, genre := range task.Genres {
		genre = strings.ToLower(genre)
		genreArrays[genre] = []types.GenreArray{}
		task.Genres[i] = genre

	}

	genreArrays["unknown"] = []types.GenreArray{}

	err, updatedTracks := pocketbase_helpers.GetInternalGenre(tracks, task.Genres, genreArrays)
	if err != "" {
		helpers.HandleError(err, user.Record.Email)
	}

	err = ai_helpers.GetExternalGenre(updatedTracks, task.Genres, genreArrays)
	if err != "" {
		helpers.HandleError(err, user.Record.Email)
		return

	}

	err, createdPlaylists := spotify_helpers.CreatePlaylists(playlistName, genreArrays)
	if err != "" {
		helpers.HandleError(err, user.Record.Email)
		return
	}

	err = pocketbase_helpers.UpdateUses(user, used, usesID)
	if err != "" {
		helpers.HandleError(err, user.Record.Email)
		return
	}

	err = pocketbase_helpers.RecordPlaylistForDeletion(user, createdPlaylists)
	if err != "" {
		helpers.HandleError(err, user.Record.Email)
		return
	}

	emailItems := []types.QueueFinishedEmailItems{}

	for i, playlist := range createdPlaylists {
		data := types.QueueFinishedEmailItems{
			ID:   i + 1,
			Name: playlist.Name,
			URL:  fmt.Sprintf("https://open.spotify.com/playlist/%s", playlist.ID),
		}

		emailItems = append(emailItems, data)
	}

	email_helpers.SendQueueFinishEmail(user.Record.Premium, used+1, emailItems, user.Record.Email)

}
