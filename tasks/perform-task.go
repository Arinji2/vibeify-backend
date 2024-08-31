package tasks

import (
	"fmt"

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

	tracks, err := spotify_helpers.GetSpotifyPlaylist(task.SpotifyURL, user)
	if err != "" {
		helpers.HandleError(err, user.Record.Email)
	}

	genreArrays := types.GenreArrays{}

	for _, genre := range task.Genres {
		genreArrays[genre] = []types.GenreArray{}

	}

	err, updatedTracks := pocketbase_helpers.GetInternalGenre(tracks, task.Genres, genreArrays)
	if err != "" {
		helpers.HandleError(err, user.Record.Email)
	}

	for _, track := range updatedTracks {
		fmt.Println(track.Track.Name)
	}

}
