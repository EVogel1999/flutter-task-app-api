package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"example.com/m/v2/database"
	"example.com/m/v2/schema"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
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
	var body schema.Category
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Error: Could not parse json."}`))
	}

	category := schema.CategoryDB{
		ID:          primitive.NewObjectID(),
		Name:        body.Name,
		Description: body.Description,
	}

	ctx := database.GetContext()
	result, err := db.Database("Task-App").Collection("Categories").InsertOne(ctx, category)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Error: Could not save to database."}`))
	}

	if result != nil {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id": "` + category.ID.Hex() + `"}`))
	}
}

func getCategory(w http.ResponseWriter, r *http.Request, db *mongo.Client) {
	params := mux.Vars(r)
	id := params["id"]

	var categoryDB schema.CategoryDB
	_id, err := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	ctx := database.GetContext()
	db.Database("Task-App").Collection("Categories").FindOne(ctx, filter).Decode(&categoryDB)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Error: Could not save to database."}`))
	}

	category := schema.Category{
		ID:          categoryDB.ID.Hex(),
		Name:        categoryDB.Name,
		Description: categoryDB.Description,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func getCategories(w http.ResponseWriter, r *http.Request, db *mongo.Client) {
	page, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
	limit := int64(10)
	skip := limit * page
	if err != nil {
		panic(err)
	}

	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.Skip = &skip

	var categories []schema.Category
	ctx := database.GetContext()
	curr, err := db.Database("Task-App").Collection("Categories").Find(ctx, bson.M{}, findOptions)
	if err != nil {
		panic(err)
	}

	for curr.Next(context.TODO()) {
		var categoryDB schema.CategoryDB
		err := curr.Decode(&categoryDB)
		if err != nil {
			panic(err)
		}

		category := schema.Category{
			ID:          categoryDB.ID.Hex(),
			Name:        categoryDB.Name,
			Description: categoryDB.Description,
		}

		categories = append(categories, category)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func editCategory(w http.ResponseWriter, r *http.Request, db *mongo.Client) {
	params := mux.Vars(r)
	id := params["id"]
	var body schema.Category
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		panic(err)
	}

	_id, err := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	ctx := database.GetContext()
	update := bson.M{
		"$set": bson.M{
			"name":        body.Name,
			"description": body.Description,
		},
	}

	db.Database("Task-App").Collection("Categories").FindOneAndUpdate(ctx, filter, update)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Category successfully updated!"}`))
}

func deleteCategory(w http.ResponseWriter, r *http.Request, db *mongo.Client) {
	params := mux.Vars(r)
	id := params["id"]
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	filter := bson.M{"_id": _id}
	ctx := database.GetContext()

	db.Database("Task-App").Collection("Categories").FindOneAndDelete(ctx, filter)
	update := bson.M{
		"$pull": bson.M{
			"category": id,
		},
	}
	db.Database("Task-App").Collection("Tasks").UpdateMany(ctx, bson.M{}, update)

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Category successfully deleted!"}`))
}
