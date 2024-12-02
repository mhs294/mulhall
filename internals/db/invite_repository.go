package db

import (
	"context"
	"fmt"
	"log"

	"github.com/mhs294/mulhall/internals/env"
	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
)

const invitesColl = "invites"

// InviteRepository is a mechanism for managing invitations for new Users to join the site.
type InviteRepository struct {
	connStr string
}

// NewInviteRepository creates a new InviteRepository instance and returns a pointer to it.
// connStr is the connection string for the MongoDB instance (e.g. - "mongodb+srv://{user}:{pass}@myinstance.mongodb.net/").
func NewInviteRepository(connStr string) (*InviteRepository, error) {
	// Create a new client and connect to the server
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, env.Timeout)
	defer cancel()
	client, err := CreateMongoDBClient(connStr, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Setup deferred connection closure for when function completes
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("%v", err)
		}
	}()

	// Test the client connection
	if err := TestMongoDBClient(client, mulhallDB, ctx); err != nil {
		return nil, fmt.Errorf("database connection test failed: %v", err)
	}

	return &InviteRepository{connStr: connStr}, nil
}

// InsertInvite inserts the provided Invite into the database.
// inv is the Invite to insert into the database.
func (ir *InviteRepository) InsertInvite(inv *types.Invite) error {
	// Create a new client and connect to the server
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, env.Timeout)
	defer cancel()
	client, err := CreateMongoDBClient(ir.connStr, ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Setup deferred connection closure for when function completes
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("%v", err)
		}
	}()

	// Acquire reference to the invites collection
	coll := client.Database(mulhallDB).Collection(invitesColl)

	// Convert the Invite into a BSON map
	bsonData, err := bson.Marshal(inv)
	if err != nil {
		return fmt.Errorf("failed to marhal invite to bson: %v", err)
	}
	var bsonMap bson.M
	err = bson.Unmarshal(bsonData, &bsonMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal invite to bson.M: %v", err)
	}

	// Insert the BSON map into the invites collection as a document
	if _, err := coll.InsertOne(ctx, bsonMap); err != nil {
		log.Fatal(err)
	}

	return nil
}

// GetInviteForEmail returns the Invite for the provided email address and token (or an error if no such Invite exists).
// email is the email address to look up the Invite for.
// token is the token string that should match with the email on the Invite.
func (ir *InviteRepository) GetInvite(email string, token string) (*types.Invite, error) {
	// Create a new client and connect to the server
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, env.Timeout)
	defer cancel()
	client, err := CreateMongoDBClient(ir.connStr, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Setup deferred connection closure for when function completes
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("%v", err)
		}
	}()

	// Load the document from the invites collection corresponding to the email address
	coll := client.Database(mulhallDB).Collection(invitesColl)
	cursor, err := coll.Find(ctx, bson.D{{Key: "email", Value: email}})
	if err != nil {
		return nil, fmt.Errorf("failed to query %s.%s: %v", mulhallDB, invitesColl, err)
	}

	// Unpack the cursor contents into a slice
	var invites []types.Invite
	if err = cursor.All(ctx, &invites); err != nil {
		return nil, fmt.Errorf("failed to parse the query results: %v", err)
	}

	// Verify that only one Invite was unpacked from the cursor
	if len(invites) == 0 {
		return nil, fmt.Errorf("no invite exists for email=%s token=%s", email, token)
	} else if len(invites) != 1 {
		return nil, fmt.Errorf("multiple invites exists for email=%s token=%s", email, token)
	}

	return &invites[0], nil
}

// AcceptInvite updates the Accepted property of the Invite keyed by the provided ID to be true
func (ir *InviteRepository) AcceptInvite(id string) error {
	// Create a new client and connect to the server
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, env.Timeout)
	defer cancel()
	client, err := CreateMongoDBClient(ir.connStr, ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Setup deferred connection closure for when function completes
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("%v", err)
		}
	}()

	// Acquire reference to the invites collection
	coll := client.Database(mulhallDB).Collection(invitesColl)

	// Define the update operation and filter
	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"accepted": true,
		},
	}

	// Perform the update
	if _, err := coll.UpdateOne(ctx, filter, update); err != nil {
		return fmt.Errorf("failed to accept invite (id=%s): %v", id, err)
	}

	return nil
}
