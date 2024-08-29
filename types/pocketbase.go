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
