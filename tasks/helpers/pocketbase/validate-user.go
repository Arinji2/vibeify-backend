package pocketbase_helpers

import (
	"encoding/json"
	"fmt"

	"github.com/Arinji2/vibeify-backend/api"
	"github.com/Arinji2/vibeify-backend/types"
)

func ValidateUser(token string) (user *types.PocketbaseUser, errorText string) {

	client := api.NewApiClient()
	res, _, err := client.SendRequestWithBody("POST", "/api/collections/users/auth-refresh", nil, map[string]string{
		"Authorization": token,
	})

	if err != nil {
		errorText = "Invalid User"
		return nil, errorText
	}
	data, err := json.Marshal(res["record"])
	if err != nil {
		fmt.Println("Marshalling Error", err)
		errorText = "Server Error"
		return nil, errorText

	}

	record := types.PocketbaseUserRecord{}

	err = json.Unmarshal(data, &record)

	if err != nil {
		fmt.Println("Error in parsing", err)
		errorText = "Server Error"
		return nil, errorText
	}

	pocketbaseUser := types.PocketbaseUser{
		Token:  token,
		Record: record,
	}

	return &pocketbaseUser, ""

}
