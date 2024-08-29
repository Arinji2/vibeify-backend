package pocketbase_helpers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Arinji2/vibeify-backend/api"
	"github.com/Arinji2/vibeify-backend/types"
)

func CheckLimit(user *types.PocketbaseUser) (used, total int, err error) {
	client := api.NewApiClient("https://db-listify.arinji.com")
	total = 0
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
		total = 0
		used = 0
		return
	}

	items, ok := res["items"].([]interface{})
	if !ok {
		total = 0
		used = 0
		return
	}

	if len(items) == 0 {
		used = 0
		return
	}

	itemMap, ok := items[0].(map[string]interface{})
	if !ok {
		fmt.Println("Item is not a map[string]interface{}")

	}

	uses, _ := strconv.Atoi(itemMap["uses"].(string))
	limit := types.PocketbaseLimit{

		Uses: uses,
	}

	used = limit.Uses

	if limit.Uses >= total {
		used = total
		if user.Record.Premium {
			err = errors.New("maximum convert requests of 10 per week reached try again next week")
		} else {
			err = errors.New("maximum convert requests of 5 per week reached please upgrade to premium to continue using the service or try again next week")
		}
	}

	return used, total, err
}
