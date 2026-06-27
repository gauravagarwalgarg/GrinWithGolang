package bootstrap

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/mongo"
)

// NewMongoDatabase creates and connects a MongoDB client.
//
// LEARNING NOTES:
// - Uses context with timeout to avoid hanging forever on connection
// - Builds URI conditionally: with auth credentials or without
// - Pings the DB to verify connectivity before returning
func NewMongoDatabase(env *Env) mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := env.DBHost
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPass

	mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)

	if dbUser == "" || dbPass == "" {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	}

	client, err := mongo.NewClient(mongodbURI)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

// CloseMongoDBConnection gracefully disconnects from MongoDB.
func CloseMongoDBConnection(client mongo.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}
