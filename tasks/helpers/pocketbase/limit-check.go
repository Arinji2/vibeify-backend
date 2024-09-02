package pocketbase_helpers

import (
	"fmt"
	"strconv"

	"github.com/Arinji2/vibeify-backend/api"
	"github.com/Arinji2/vibeify-backend/types"
)

func CheckLimit(user *types.PocketbaseUser) (used int, usesID string, errorText string) {
	client := api.NewApiClient("https://db-listify.arinji.com")
	total := 0
	used = 0
	if user.Record.Premium {
		total = 10
	} else {
		total = 5
	}
	res, _, err := client.SendRequestWithQuery("GET", "/api/collections/convertLimit/records", map[string]string{
		"page":    "1",
		"perPage": "1",
	}, map[string]string{
		"Authorization": user.Token,
	})

	if err != nil {

		return
	}

	items, ok := res["items"].([]interface{})
	if !ok {

		return
	}

	if len(items) == 0 {

		return
	}

	itemMap, ok := items[0].(map[string]interface{})
	if !ok {
		fmt.Println("Item is not a map[string]interface{}")

	}

	uses, _ := strconv.Atoi(itemMap["uses"].(string))
	usesID = itemMap["id"].(string)
	limit := types.PocketbaseLimit{

		Uses: uses,
	}

	used = limit.Uses

	if limit.Uses >= total {
		used = total
		if user.Record.Premium {
			errorText = "Maximum convert requests of 10 per week reached try again next week"
		} else {
			errorText = "Maximum convert requests of 5 per week reached please upgrade to premium to continue using the service or try again next week"
		}
	}

	return used, usesID, errorText
}
