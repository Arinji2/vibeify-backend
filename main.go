package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Arinji2/vibeify-backend/tasks"
	indexing_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/pocketbase/indexing"
	"github.com/Arinji2/vibeify-backend/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/gookit/slog"
	"github.com/joho/godotenv"
)

var tasksArr []types.AddTaskType
var (
	taskInProgress bool = false
	mu             sync.Mutex
)

func SkipLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}
		middleware.Logger(next).ServeHTTP(w, r)
	})
}

func main() {

	r := chi.NewRouter()
	r.Use(SkipLoggingMiddleware)

	err := godotenv.Load()

	if err != nil {
		isProduction := os.Getenv("ENVIRONMENT") == "PRODUCTION"
		if !isProduction {
			log.Fatal("Error loading .env file")
		} else {
			slog.Warn("Using Production Environment")
		}
	} else {
		slog.Warn("Using Development Environment")
	}

	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Vibeify Backend: Request Received")
		defer w.Write([]byte("Vibeify Backend: Request Received"))
		key := r.URL.Query()["key"]

		if len(key) != 0 {

			if key[0] != os.Getenv("ACCESS_KEY") {
				render.Status(r, 401)
				return

			}

		}

		render.Status(r, 200)
	})

	r.Post("/addTask", func(w http.ResponseWriter, r *http.Request) {

		var requestBody types.AddTaskType

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&requestBody)
		if err != nil {
			http.Error(w, "Invalid Input", http.StatusBadRequest)
			return
		}

		tasksArr = append(tasksArr, requestBody)

		render.Status(r, 200)
	})

	go checkTasks()
	go checkIndexing()
	go checkPlaylistIndexing()

	r.Get("/index/playlists", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query()["key"]

		if len(key) != 0 {

			if key[0] != os.Getenv("ACCESS_KEY") {
				render.Status(r, 401)
				return

			}

		}

		go indexing_helpers.CheckPlaylistIndexing()
		render.Status(r, 200)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("Vibeify Backend: Health Check"))
		render.Status(r, 200)
	})

	http.ListenAndServe(":8080", r)

}

func checkTasks() {
	mu.Lock()
	defer mu.Unlock()

	if taskInProgress {
		return
	}
	ticker := time.NewTicker(time.Millisecond * 10)
	slog.Info("Cron Job For Task Check Started")
	for range ticker.C {
		if len(tasksArr) > 0 {
			selectedTask := tasksArr[0]
			slog.Info(fmt.Sprintf("New task found. Tasks Remaining: %v", len(tasksArr)))
			taskInProgress = true
			tasksArr = tasksArr[1:]
			tasks.PerformTask(selectedTask)
			taskInProgress = false

		}
	}
}

func checkIndexing() {
	ticker := time.NewTicker(time.Second * 10)
	slog.Info("Cron Job For Indexing Started")
	for range ticker.C {
		indexing_helpers.CheckIndexing()
		indexing_helpers.CleanupIndexing()
	}
}

func checkPlaylistIndexing() {

	ticker := time.NewTicker(time.Hour * 24)
	slog.Info("Cron Job For Playlist Indexing Started")
	for range ticker.C {
		indexing_helpers.CheckPlaylistIndexing()
		indexing_helpers.CleanupIndexing()
	}
}
