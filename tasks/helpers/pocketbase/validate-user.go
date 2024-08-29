package pocketbase_helpers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Arinji2/vibeify-backend/api"
	"github.com/Arinji2/vibeify-backend/types"
)

func ValidateUser(token string) (user *types.PocketbaseUser, err error) {

	client := api.NewApiClient("https://db-listify.arinji.com")
	res, _, err := client.SendRequestWithBody("POST", "/api/collections/users/auth-refresh", nil, map[string]string{
		"Authorization": token,
	})

	if err != nil {
		return nil, errors.New("invalid user")
	}
	data, err := json.Marshal(res["record"])
	if err != nil {
		fmt.Println("Marshalling Error", err)
		return nil, errors.New("server error")

	}

	record := types.PocketbaseUserRecord{}

	err = json.Unmarshal(data, &record)

	if err != nil {
		fmt.Println("Error in parsing", err)
		return nil, errors.New("server error")
	}

	pocketbaseUser := types.PocketbaseUser{
		Token:  token,
		Record: record,
	}

	return &pocketbaseUser, nil

}
