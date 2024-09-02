package ai_helpers

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Arinji2/vibeify-backend/api"
	"github.com/Arinji2/vibeify-backend/types"
	"github.com/joho/godotenv"
)

func GetExternalGenre(remainingTracks []types.SpotifyPlaylistItem, genres []string, genreArrays types.GenreArrays, updatedPrompt ...string) (errorString string) {
	godotenv.Load()
	accessKey := os.Getenv("ACCESS_KEY")

	client := api.NewApiClient("https://ai.arinji.com")

	var wg sync.WaitGroup
	var mu sync.Mutex
	var pool = make(chan struct{}, 10)

	for _, track := range remainingTracks {
		wg.Add(1)
		pool <- struct{}{}
		go func(track types.SpotifyPlaylistItem) {
			defer wg.Done()
			defer func() {
				<-pool
			}()

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

			retries := 0
			maxRetries := 3
			for retries < maxRetries {
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
					if status == 500 {
						//this is when the AI API is overloaded, we wait here
						time.Sleep(time.Minute * 1)
					}
					continue
				}

				message, ok := res["message"].(string)
				if !ok {
					fmt.Println("Error converting message to string")
					retries++
					continue
				}

				message = strings.ToLower(message)

				if slices.Contains(genres, message) {
					mu.Lock()
					genreArrays[message] = append(genreArrays[message], types.GenreArray{
						URI: track.Track.URI,
					})
					mu.Unlock()
					break
				} else {
					updatedPromptString := fmt.Sprintf("The genre you guessed doesn't exist. DON'T GUESS THE GENRE %s", message)
					prompt = prompt + updatedPromptString
					retries++
					fmt.Println("Retrying AI Request With Updated Prompt")
				}
			}

			if retries == maxRetries {
				mu.Lock()
				genreArrays["unknown"] = append(genreArrays["unknown"], types.GenreArray{
					URI: track.Track.URI,
				})
				mu.Unlock()
			}

		}(track)
	}

	wg.Wait()
	errorString = ""

	return
}
