package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func setTaskHandler(router *mux.Router, db *mongo.Client) {
	router.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		createTask(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		getTask(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		getTasks(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		editTask(w, r, db)
	}).Methods("PUT")

	router.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteTask(w, r, db)
	}).Methods("DELETE")
}

func createTask(w http.ResponseWriter, r *http.Request, db *mongo.Client) {

}

func getTask(w http.ResponseWriter, r *http.Request, db *mongo.Client) {

}

func getTasks(w http.ResponseWriter, r *http.Request, db *mongo.Client) {

}

func editTask(w http.ResponseWriter, r *http.Request, db *mongo.Client) {

}

func deleteTask(w http.ResponseWriter, r *http.Request, db *mongo.Client) {

}
