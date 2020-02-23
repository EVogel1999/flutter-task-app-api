package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
}

type TaskDB struct {
	ID          primitive.ObjectID `json:"_id"`
	Name        string             `json:"name"`
	Date        time.Time          `json:"date"`
	Description string             `json:"description"`
	Category    string             `json:"category"`
}
