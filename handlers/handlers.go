package handlers

import (
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUpHandlers(router *mux.Router, client *mongo.Client) {
	setTaskHandler(router, client)
	setCategoryHandler(router, client)
}
