package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Arinji2/vibeify-backend/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	err := godotenv.Load()
	if err != nil {
		isProduction := os.Getenv("ENVIRONMENT") == "PRODUCTION"
		if !isProduction {
			log.Fatal("Error loading .env file")
		} else {
			fmt.Println("Using Production Environment")
		}
	} else {
		fmt.Println("Using Development Environment")
	}

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

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		render.Status(r, 200)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Vibeify Backend: Health Check")
		w.Write([]byte("Vibeify Backend: Health Check"))
		render.Status(r, 200)
	})

	http.ListenAndServe(":8080", r)
}
