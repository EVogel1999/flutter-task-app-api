package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func setCategoryHandler(router *mux.Router, db *mongo.Client) {
	router.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		createCategory(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		getCategory(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		getCategories(w, r, db)
	}).Methods("GET").Queries("page", "{[0-9]}")

	router.HandleFunc("/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		editCategory(w, r, db)
	}).Methods("PUT")

	router.HandleFunc("/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteCategory(w, r, db)
	}).Methods("DELETE")
}

func createCategory(w http.ResponseWriter, r *http.Request, db *mongo.Client) {

}

func getCategory(w http.ResponseWriter, r *http.Request, db *mongo.Client) {

}

func getCategories(w http.ResponseWriter, r *http.Request, db *mongo.Client) {

}

func editCategory(w http.ResponseWriter, r *http.Request, db *mongo.Client) {

}

func deleteCategory(w http.ResponseWriter, r *http.Request, db *mongo.Client) {

}
