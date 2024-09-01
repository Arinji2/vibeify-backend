package ai_helpers

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Arinji2/vibeify-backend/api"
	"github.com/Arinji2/vibeify-backend/types"
	"github.com/joho/godotenv"
)

func GetExternalGenre(remainingTracks []types.SpotifyPlaylistItem, genres []string, genreArrays types.GenreArrays, updatedPrompt ...string) (errorString string) {
	godotenv.Load()
	accessKey := os.Getenv("ACCESS_KEY")

	client := api.NewApiClient("https://ai.arinji.com")

	retries := 0

	for _, track := range remainingTracks {

		artistString := ""
		genreString := ""
		for _, artist := range track.Track.Artists {
			artistString = artistString + artist.Name + ", "
		}
		for _, genre := range genres {
			genreString = genreString + genre + ", "
		}

		prompt := fmt.Sprintf("Given the song name %s by artists %s Your objective is to guess the genre of the song, which is ONLY %s. Reply with ONLY the genre name, nothing else. Choose only one genre", track.Track.Name, artistString, genreString)
		if len(updatedPrompt) > 0 {
			prompt = prompt + updatedPrompt[0]
		}
		body := []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		}
		headers := map[string]string{
			"Content-Type":  "application/json",
			"Authorization": accessKey,
		}

		res, status, err := client.SendRequestWithBody("POST", "/completions", body, headers)
		if err != nil || status != 200 {
			retries++
			if retries < 3 {
				fmt.Println("Retrying AI Request")
				GetExternalGenre(remainingTracks, genres, genreArrays)
			} else {
				genreArrays["unknown"] = append(genreArrays["unknown"], types.GenreArray{
					URI:  track.Track.URI,
					Name: track.Track.Name,
				})
			}

		}

		message, ok := res["message"].(string)
		if !ok {
			fmt.Println("Error converting message to string")
			genreArrays["unknown"] = append(genreArrays["unknown"], types.GenreArray{
				URI:  track.Track.URI,
				Name: track.Track.Name,
			})

		}

		message = strings.ToLower(message)

		if !slices.Contains(genres, message) {
			updatedPromptString := fmt.Sprintf("The genre you guessed dosent exist. DONT GUESS THAT GENRE %s", message)
			retries++
			if retries < 3 {
				fmt.Println("Retrying AI Request With Updated Prompt")
				GetExternalGenre(remainingTracks, genres, genreArrays, updatedPromptString)
			} else {
				genreArrays["unknown"] = append(genreArrays["unknown"], types.GenreArray{
					URI:  track.Track.URI,
					Name: track.Track.Name,
				})

			}
		}

		for genre := range genreArrays {
			if genre == message {
				genreArrays[genre] = append(genreArrays[genre], types.GenreArray{
					URI:  track.Track.URI,
					Name: track.Track.Name,
				})

			}
		}

	}

	errorString = ""

	return
}
