package types

type AddTaskType struct {
	SpotifyURL string   `json:"spotifyURL"`
	UserToken  string   `json:"userToken"`
	Genres     []string `json:"genres"`
}
