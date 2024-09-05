package pocketbase_helpers

import (
	"errors"
	"strconv"

	"github.com/Arinji2/vibeify-backend/api"
	"github.com/Arinji2/vibeify-backend/types"
	user_errors "github.com/Arinji2/vibeify-backend/user-errors"
)

func CheckLimit(user *types.PocketbaseUser) (used int, usesID string, err error) {
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

		return used, usesID, user_errors.NewUserError("", err)
	}

	items, ok := res["items"].([]interface{})
	if !ok {

		return used, usesID, user_errors.NewUserError("", err)
	}

	if len(items) == 0 {

		return used, usesID, nil
	}

	itemMap, ok := items[0].(map[string]interface{})
	if !ok {

		return used, usesID, user_errors.NewUserError("", errors.New("item is not a map[string]interface{}"))
	}

	uses, _ := strconv.Atoi(itemMap["uses"].(string))
	usesID = itemMap["id"].(string)

	limit := types.PocketbaseLimit{

		Uses: uses,
	}

	used = limit.Uses

	if limit.Uses >= total {

		if user.Record.Premium {
			return uses, usesID, user_errors.NewUserError("Maximum convert requests of 10 per week reached please upgrade to premium to continue using the service or try again next week", errors.New("max paid limit reached"))
		} else {
			return uses, usesID, user_errors.NewUserError("Maximum convert requests of 5 per week reached please upgrade to premium to continue using the service or try again next week", errors.New("max free limit reached"))
		}
	}

	return used, usesID, nil
}
