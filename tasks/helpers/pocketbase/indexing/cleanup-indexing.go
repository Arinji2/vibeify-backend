package indexing_helpers

import (
	"fmt"
	"math"
	"sync"

	"github.com/Arinji2/vibeify-backend/api"
	pocketbase_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/pocketbase"
)

const MaxRecordsAllowed = 5000

func CleanupIndexing() {
	adminToken, err := pocketbase_helpers.GetPocketbaseAdminToken()

	if err != "" {
		fmt.Println(err)
		return
	}
	client := api.NewApiClient()
	res, _, error := client.SendRequestWithQuery("GET", "/api/collections/songs/records", map[string]string{
		"page":    "1",
		"perPage": "1",
		"fields":  "id",
	}, map[string]string{
		"Authorization": adminToken,
	})

	if error != nil {
		fmt.Println(error)
		return
	}

	totalRecords, ok := res["totalItems"].(float64)
	if !ok {
		fmt.Println("Error getting total records")
	}

	if totalRecords < MaxRecordsAllowed {
		return
	}

	perPage := math.Max((totalRecords - MaxRecordsAllowed), MaxRecordsAllowed)

	res, _, error = client.SendRequestWithQuery("GET", "/api/collections/songs/records", map[string]string{
		"page":      "1",
		"perPage":   fmt.Sprintf("%v", perPage),
		"fields":    "id",
		"sort":      "totalUses",
		"skipTotal": "true",
	}, map[string]string{
		"Authorization": adminToken,
	})

	if error != nil {
		fmt.Println(error)
		return
	}

	items, ok := res["items"].([]interface{})
	if !ok {
		fmt.Println("Error getting items")
		return
	}

	pool := make(chan struct{}, 5)
	cleanupWg := sync.WaitGroup{}

	fmt.Println("Extra Results Found, Purging.")

	for _, item := range items {
		cleanupWg.Add(1)
		pool <- struct{}{}
		go func(item interface{}) {
			defer cleanupWg.Done()
			defer func() { <-pool }()

			id, ok := item.(map[string]interface{})["id"]
			if !ok {
				fmt.Println("Error getting id")
				return
			}

			_, _, error := client.SendRequestWithQuery("DELETE", "/api/collections/songs/records/"+id.(string), nil, map[string]string{
				"Authorization": adminToken,
			})

			if error != nil {
				fmt.Println(error)
				return
			}

		}(item)
	}
}
