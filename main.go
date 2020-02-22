package main

import (
	"fmt"
	"net/http"
	"os"

	"example.com/m/v2/database"
	"example.com/m/v2/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func loadENVs() {
	if err := godotenv.Load(); err != nil {
		fmt.Print("\nNo .env file found")
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Welcome to the Flutter Task App API!"}`))
}

func main() {
	loadENVs()

	client := database.Connect()
	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler)
	handlers.SetUpHandlers(r, client)

	port, exists := os.LookupEnv("PORT")

	if exists {
		fmt.Printf("\nListening on PORT " + port)
		http.ListenAndServe(":"+port, r)
	} else {
		fmt.Printf("\nListening on PORT 3000")
		http.ListenAndServe(":3000", r)
	}
}
