package pocketbase_helpers

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Arinji2/vibeify-backend/api"
)

var (
	tokenCache  string
	expiryCache time.Time
	mu          sync.Mutex
)

const tokenValidity = 604800 * time.Second // 7 days
func GetPocketbaseAdminToken() (token string, errorString string) {

	mu.Lock()
	defer mu.Unlock()
	errorString = "Server Error"
	if time.Now().Before(expiryCache) && tokenCache != "" {
		token = tokenCache
		errorString = ""
		return
	}

	identityEmail := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")

	if identityEmail == "" || password == "" {
		fmt.Println("Environment Variables not present to authenticate Admin")
		return

	}

	body := map[string]string{
		"identity": identityEmail,
		"password": password,
	}

	client := api.NewApiClient()
	result, _, err := client.SendRequestWithBody("POST", "/api/admins/auth-with-password", body, nil)

	if err != nil {
		fmt.Println("Admin Login failed:", err)
		return
	}

	token, ok := result["token"].(string)
	if !ok || token == "" {
		fmt.Println("Admin Token not found or not a string")
		return
	}

	tokenCache = token
	expiryCache = time.Now().Add(tokenValidity)

	errorString = ""
	return

}
