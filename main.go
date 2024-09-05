package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	custom_log "github.com/Arinji2/vibeify-backend/logger"
	"github.com/Arinji2/vibeify-backend/tasks"
	pocketbase_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/pocketbase"
	indexing_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/pocketbase/indexing"
	"github.com/Arinji2/vibeify-backend/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
)

type TaskManager struct {
	tasks          []types.AddTaskType
	taskInProgress bool
	mu             sync.Mutex
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
			custom_log.Logger.Warn("Using Production Environment")
		}
	} else {
		custom_log.Logger.Warn("Using Development Environment")
	}

	taskManager := &TaskManager{}
	go taskManager.startTaskWorker()
	go startCronJobs()

	r.Get("/", healthHandler)
	r.Post("/addTask", taskManager.addTaskHandler)
	r.Get("/index/playlists", playlistIndexingHandler)
	r.Get("/health", healthCheckHandler)

	http.ListenAndServe(":8080", r)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Vibeify Backend: Request Received")
	w.Write([]byte("Vibeify Backend: Request Received"))
	key := r.URL.Query().Get("key")

	if key != "" && key != os.Getenv("ACCESS_KEY") {
		render.Status(r, http.StatusUnauthorized)
		return
	}

	render.Status(r, http.StatusOK)
}

func (tm *TaskManager) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody types.AddTaskType
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, task := range tm.tasks {
		if task.UserToken == requestBody.UserToken {
			http.Error(w, "Task Already Exists", http.StatusBadRequest)
			return
		}
	}
	tm.tasks = append(tm.tasks, requestBody)

	render.Status(r, http.StatusOK)
}

func playlistIndexingHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key != "" && key != os.Getenv("ACCESS_KEY") {
		render.Status(r, http.StatusUnauthorized)
		return
	}

	go indexing_helpers.CheckPlaylistIndexing()
	render.Status(r, http.StatusOK)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Vibeify Backend: Health Check"))
	render.Status(r, http.StatusOK)
}

func (tm *TaskManager) startTaskWorker() {
	ticker := time.NewTicker(10 * time.Millisecond)
	custom_log.Logger.Info("Cron Job For Task Check Started")

	for range ticker.C {
		tm.mu.Lock()
		if tm.taskInProgress || len(tm.tasks) == 0 {
			tm.mu.Unlock()
			continue
		}

		selectedTask := tm.tasks[0]
		tm.tasks = tm.tasks[1:]
		tm.taskInProgress = true
		tm.mu.Unlock()

		custom_log.Logger.Info(fmt.Sprintf("New task found. Tasks Remaining: %v", len(tm.tasks)))
		tasks.PerformTask(selectedTask)

		tm.mu.Lock()
		tm.taskInProgress = false
		tm.mu.Unlock()
	}
}

func startIndexingJobs() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		custom_log.Logger.Info("Cron Job For Priority Indexing Started")
		for range ticker.C {
			indexing_helpers.CheckIndexing()
			indexing_helpers.CleanupIndexing()
		}
	}()

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		custom_log.Logger.Info("Cron Job For Playlist Indexing Started")
		for range ticker.C {
			indexing_helpers.CheckPlaylistIndexing()
			indexing_helpers.CleanupIndexing()
		}
	}()
}

func startCronJobs() {
	go startIndexingJobs()

	ticker := time.NewTicker(24 * time.Hour)
	custom_log.Logger.Info("Cron Job For Playlist Deletion Started")
	for range ticker.C {
		go pocketbase_helpers.DeleteExpiredPlaylists()
	}

}

func SkipLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}
		middleware.Logger(next).ServeHTTP(w, r)
	})
}
