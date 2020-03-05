package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryDB struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
}
