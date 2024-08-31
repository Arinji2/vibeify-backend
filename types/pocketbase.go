package types

type PocketbaseUser struct {
	Token  string               `json:"token"`
	Record PocketbaseUserRecord `json:"record"`
}

type PocketbaseUserRecord struct {
	ID              string `json:"id"`
	CollectionID    string `json:"collectionId"`
	CollectionName  string `json:"collectionName"`
	Username        string `json:"username"`
	Verified        bool   `json:"verified"`
	EmailVisibility bool   `json:"emailVisibility"`
	Email           string `json:"email"`
	Created         string `json:"created"`
	Updated         string `json:"updated"`
	Name            string `json:"name"`
	Dicebear        string `json:"dicebear"`
	Premium         bool   `json:"premium"`
}

type PocketbaseLimit struct {
	CollectionId   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	Created        string `json:"created"`
	Id             string `json:"id"`
	Updated        string `json:"updated"`
	User           string `json:"user"`
	Uses           int    `json:"uses"`
}

type PocketbaseSongRecord struct {
	ID             string   `json:"id"`
	CollectionID   string   `json:"collectionId"`
	CollectionName string   `json:"collectionName"`
	Created        string   `json:"created"`
	Updated        string   `json:"updated"`
	Genres         []string `json:"genres"`
	SpotifyID      string   `json:"spotifyID"`
	Name           string   `json:"name"`
}
