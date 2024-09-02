package types

type GenreArrays map[string][]GenreArray

type GenreArray struct {
	URI string `json:"uri"`
}
