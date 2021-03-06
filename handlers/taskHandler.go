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

func setTaskHandler(router *mux.Router, db *mongo.Client) {
	router.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		createTask(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		getTask(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		getTasks(w, r, db)
	}).Methods("GET").Queries("page", "{[0-9]}")

	router.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		editTask(w, r, db)
	}).Methods("PUT")

	router.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteTask(w, r, db)
	}).Methods("DELETE")
}

func createTask(w http.ResponseWriter, r *http.Request, db *mongo.Client) {
	var body schema.Task
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Error: Could not parse json."}`))
	}

	task := schema.TaskDB{
		ID:          primitive.NewObjectID(),
		Name:        body.Name,
		Description: body.Description,
		Date:        body.Date,
		Category:    body.Category,
	}

	ctx := database.GetContext()
	result, err := db.Database("Task-App").Collection("Tasks").InsertOne(ctx, task)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Error: Could not save to database."}`))
	}

	if result != nil {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id": "` + task.ID.Hex() + `"}`))
	}
}

func getTask(w http.ResponseWriter, r *http.Request, db *mongo.Client) {
	params := mux.Vars(r)
	id := params["id"]

	var taskDB schema.TaskDB
	_id, err := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	ctx := database.GetContext()
	db.Database("Task-App").Collection("Tasks").FindOne(ctx, filter).Decode(&taskDB)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Error: Could not save to database."}`))
	}

	task := schema.Task{
		ID:          taskDB.ID.Hex(),
		Name:        taskDB.Name,
		Description: taskDB.Description,
		Date:        taskDB.Date,
		Category:    taskDB.Category,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func getTasks(w http.ResponseWriter, r *http.Request, db *mongo.Client) {
	page, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
	limit := int64(10)
	skip := limit * page
	if err != nil {
		panic(err)
	}

	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.Skip = &skip

	var tasks []schema.Task
	ctx := database.GetContext()
	curr, err := db.Database("Task-App").Collection("Tasks").Find(ctx, bson.M{}, findOptions)
	if err != nil {
		panic(err)
	}

	for curr.Next(context.TODO()) {
		var taskDB schema.TaskDB
		err := curr.Decode(&taskDB)
		if err != nil {
			panic(err)
		}

		task := schema.Task{
			ID:          taskDB.ID.Hex(),
			Name:        taskDB.Name,
			Description: taskDB.Description,
			Date:        taskDB.Date,
			Category:    taskDB.Category,
		}

		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func editTask(w http.ResponseWriter, r *http.Request, db *mongo.Client) {
	params := mux.Vars(r)
	id := params["id"]
	var body schema.Task
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
			"category":    body.Category,
			"description": body.Description,
			"date":        body.Date,
		},
	}

	db.Database("Task-App").Collection("Tasks").FindOneAndUpdate(ctx, filter, update)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Task successfully updated!"}`))
}

func deleteTask(w http.ResponseWriter, r *http.Request, db *mongo.Client) {
	params := mux.Vars(r)
	id := params["id"]
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	filter := bson.M{"_id": _id}
	ctx := database.GetContext()

	db.Database("Task-App").Collection("Tasks").FindOneAndDelete(ctx, filter)

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Task successfully deleted!"}`))
}
