package tasks

import (
	"fmt"
	"strings"

	"github.com/Arinji2/vibeify-backend/tasks/helpers"
	ai_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/ai"
	email_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/emails"
	pocketbase_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/pocketbase"
	indexing_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/pocketbase/indexing"
	spotify_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/spotify"
	"github.com/Arinji2/vibeify-backend/types"
)

func PerformTask(task types.AddTaskType) {
	user, err := pocketbase_helpers.ValidateUser(task.UserToken)
	if err != nil {
		helpers.HandleError(err, "")
	}

	used, usesID, err := pocketbase_helpers.CheckLimit(user)

	if err != nil {
		helpers.HandleError(err, user.Record.Email)
	}

	email_helpers.SendQueueAdditionEmail(user.Record.Premium, user.Record.Email)

	tracks, playlistName, err := spotify_helpers.GetSpotifyPlaylist(task.SpotifyURL, user)
	if err != nil {
		helpers.HandleError(err, user.Record.Email)
	}

	genreArrays := types.GenreArrays{}

	for i, genre := range task.Genres {
		genre = strings.ToLower(genre)
		genreArrays[genre] = []types.GenreArray{}
		task.Genres[i] = genre

	}

	genreArrays["unknown"] = []types.GenreArray{}

	updatedTracks, err := pocketbase_helpers.GetInternalGenre(tracks, task.Genres, genreArrays)
	if err != nil {
		helpers.HandleError(err, user.Record.Email)
	}

	ai_helpers.GetExternalGenre(updatedTracks, task.Genres, genreArrays)

	createdPlaylists, err := spotify_helpers.CreatePlaylists(playlistName, genreArrays)
	if err != nil {
		helpers.HandleError(err, user.Record.Email)
	}

	err = pocketbase_helpers.UpdateUses(user, used, usesID)
	if err != nil {
		helpers.HandleError(err, user.Record.Email)
	}

	err = pocketbase_helpers.RecordPlaylistForDeletion(user, createdPlaylists)
	if err != nil {
		helpers.HandleError(err, user.Record.Email)
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

	go indexing_helpers.QueueSongIndexing(updatedTracks, "1")

}
