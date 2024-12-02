package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mulhallDB = "mulhall"

// CreateMongoDBClient creates a client for a MongoDB instance using the provided connection string and Context.
// connStr is the connection string for the MongoDB instance (e.g. - "mongodb+srv://{user}:{pass}@myinstance.mongodb.net/").
// ctx is the context.Context within which the MongoDB connection will be established.
func CreateMongoDBClient(connStr string, ctx context.Context) (*mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connStr).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// TestMongoDBClient executes a quick test to validate the connection for an established MongoDB client.
// c is the pointer to the MongoDB client instance.
// db is the name of the database to test the client against.
// ctx is the context.Context within which the client test should be executed.
func TestMongoDBClient(c *mongo.Client, db string, ctx context.Context) error {
	var result bson.M
	if err := c.Database(db).RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return err
	}

	return nil
}
