package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Client {
	fmt.Printf("\nConnecting to database...")

	connection, connectionExists := os.LookupEnv("DB_CONNECTION")
	dbPort, portExists := os.LookupEnv("DB_PORT")

	if connectionExists && portExists {
		connection = strings.ReplaceAll(connection, "<DB_PORT>", dbPort)
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(connection)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nConnected to database!")

	return client
}
