package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB represents a MongoDB instance and abstracts the various operations it provides.
type MongoDB struct {
	connStr string
	timeout time.Duration
}

// NewMongoDB creates a new instance of a MongoDB and returns a pointer to it.
//
// connStr is the connection string for the MongoDB instance (e.g. - "mongodb+srv://{user}:{pass}@myinstance.mongodb.net/").
//
// timeout is the [time.Duration] specifying the timeout for any database operations performed with the MongoDB instance.
func NewMongoDB(connStr string, timeout time.Duration) *MongoDB {
	return &MongoDB{connStr: connStr, timeout: timeout}
}

// TestConnection tests the connection to the MongoDB instance and returns any error that occurs.
//
// dbName is the name of the database to test the connection against.
func (mdb *MongoDB) TestConnection(dbName string) error {
	// Create a new client and connect to the server
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, mdb.timeout)
	defer cancel()
	client, err := createClient(mdb.connStr, ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Setup deferred connection closure for when function completes
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("%v", err)
		}
	}()

	// Test the client connection
	var result bson.M
	if err := client.Database(dbName).RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return fmt.Errorf("database connection test failed: %v", err)
	}

	return nil
}

// GetAll loads all documents from the specified database and collection using the specified query into the provided results object.
//
// dbName is the name of the database to query.
//
// collName is the name of the collection to query.
//
// query is the bson.D representing the query to load the desired documents.
//
// results is the provided object into which the documents returned from the query will be deserialized and stored.
func (mdb *MongoDB) GetAll(dbName string, collName string, query bson.D, results any) error {
	// Create a new client and connect to the server
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, mdb.timeout)
	defer cancel()
	client, err := createClient(mdb.connStr, ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Setup deferred connection closure for when function completes
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("%v", err)
		}
	}()

	// Load the documents from the database collection using the specified query filters
	coll := client.Database(dbName).Collection(collName)
	cursor, err := coll.Find(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to query %s.%s: %v", dbName, collName, err)
	}

	// Unpack the cursor contents into results parameter
	if err = cursor.All(ctx, results); err != nil {
		return fmt.Errorf("failed to parse the query results: %v", err)
	}

	return nil
}

// GetOne loads the first document from the specified database and collection matching the specified query into the provided result object.
//
// dbName is the name of the database to query.
//
// collName is the name of the collection to query.
//
// query is the bson.D representing the query to load the desired document.
//
// result is the provided object into which the document returned from the query will be deserialized and stored.
func (mdb *MongoDB) GetOne(dbName string, collName string, query bson.D, result any) error {
	// Create a new client and connect to the server
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, mdb.timeout)
	defer cancel()
	client, err := createClient(mdb.connStr, ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Setup deferred connection closure for when function completes
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("%v", err)
		}
	}()

	// Load the document from the database collection using the specified query filters
	coll := client.Database(dbName).Collection(collName)
	qRes := coll.FindOne(ctx, query)
	if err = qRes.Err(); err == mongo.ErrNoDocuments {
		// Couldn't find any document for the query, defer to caller
		return nil
	} else if err != nil {
		// An unexpected error occurred, return it to the caller
		return fmt.Errorf("failed to query %s.%s: %v", dbName, collName, err)
	}

	// Deserialize the returned document into the result
	if err = qRes.Decode(result); err != nil {
		return fmt.Errorf("failed to parse the query result: %v", err)
	}

	return nil
}

// InsertOne inserts the provided object as a document into the specified database collection.
//
// dbName is the name of the database where the document will be inserted.
//
// collName is the name of the collection where the document will be inserted.
//
// doc is the object that will be serialized into a document and inserted into the database collection.
func (mdb *MongoDB) InsertOne(dbName string, collName string, doc any) error {
	// Create a new client and connect to the server
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, mdb.timeout)
	defer cancel()
	client, err := createClient(mdb.connStr, ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Setup deferred connection closure for when function completes
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("%v", err)
		}
	}()

	// Acquire reference to the database collection
	coll := client.Database(dbName).Collection(collName)

	// Convert the document into a BSON map
	bsonData, err := bson.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marhal document to bson: %v", err)
	}
	var bsonMap bson.M
	err = bson.Unmarshal(bsonData, &bsonMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal document to bson.M: %v", err)
	}

	// Insert the BSON map into the database collection
	if _, err := coll.InsertOne(ctx, bsonMap); err != nil {
		log.Fatal(err)
	}

	return nil
}

// UpdateOne updates a single document in the specified database collection using the provided filter and update queries.
//
// dbName is the name of the database containing the document to update.
//
// collName is the name of the collection containing the document to update.
//
// filter is the bson.M representing the query to find the document to update.
//
// update is the bson.M representing the query to update the found document.
func (mdb *MongoDB) UpdateOne(dbName string, collName string, filter bson.M, update bson.M) error {
	// Create a new client and connect to the server
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, mdb.timeout)
	defer cancel()
	client, err := createClient(mdb.connStr, ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Setup deferred connection closure for when function completes
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("%v", err)
		}
	}()

	// Acquire reference to the database collection
	coll := client.Database(dbName).Collection(collName)

	// Perform the update
	if _, err := coll.UpdateOne(ctx, filter, update); err != nil {
		return fmt.Errorf("failed to update document (filter=%v, update=%v): %v", filter, update, err)
	}

	return nil
}

func createClient(connStr string, ctx context.Context) (*mongo.Client, error) {
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
