package pocketbase_helpers

import (
	"fmt"

	"github.com/Arinji2/vibeify-backend/api"
	"github.com/Arinji2/vibeify-backend/types"
)

func UpdateUses(user *types.PocketbaseUser, used int, usesID string) (errorText string) {
	adminToken, err := GetPocketbaseAdminToken()
	errorText = "Server Error"
	if err != "" {
		fmt.Println(err)
		return
	}

	client := api.NewApiClient()
	_, _, error := client.SendRequestWithBody("POST", fmt.Sprintf("/api/collections/convertLimit/records/%s", usesID), map[string]string{
		"uses": fmt.Sprintf("%d", used+1),
	}, map[string]string{
		"Authorization": adminToken,
	})

	if error != nil {
		fmt.Println(error)
		return
	}

	errorText = ""

	return
}
