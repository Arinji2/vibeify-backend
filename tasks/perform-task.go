package tasks

import (
	"fmt"
	"strings"

	"github.com/Arinji2/vibeify-backend/tasks/helpers"
	pocketbase_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/pocketbase"
	spotify_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/spotify"
	"github.com/Arinji2/vibeify-backend/types"
)

func PerformTask(task types.AddTaskType) {
	user, err := pocketbase_helpers.ValidateUser(task.UserToken)
	if err != "" {
		helpers.HandleError(err, "")
	}

	used, total, err := pocketbase_helpers.CheckLimit(user)
	fmt.Println(used, total)
	if err != "" {
		helpers.HandleError(err, user.Record.Email)
	}

	//email_helpers.SendQueueAdditionEmail(user.Record.Premium, user.Record.Email)
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

	err, _ = pocketbase_helpers.GetInternalGenre(tracks, task.Genres, genreArrays)
	if err != "" {
		helpers.HandleError(err, user.Record.Email)
	}

	// err = ai_helpers.GetExternalGenre(updatedTracks, task.Genres, genreArrays)
	// if err != "" {
	// 	helpers.HandleError(err, user.Record.Email)
	// 	return

	// }

	err, createdPlaylists := spotify_helpers.CreatePlaylists(playlistName, genreArrays)
	if err != "" {
		helpers.HandleError(err, user.Record.Email)
		return
	}
	fmt.Println(createdPlaylists)

}
