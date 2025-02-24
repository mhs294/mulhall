package repos

import (
	"fmt"

	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
)

// InviteRepository manages Invite records in the database.
type InviteRepository struct {
	mdb      *db.MongoDB
	dbName   string
	collName string
}

// NewInviteRepository creates a new InviteRepository instance and returns a pointer to it.
//
// db is the MongoDB instance used by the InviteRepository.
func NewInviteRepository(mdb *db.MongoDB) *InviteRepository {
	return &InviteRepository{mdb: mdb, dbName: "mulhall", collName: "invites"}
}

// TestConnection tests the connection to the MongoDB instance and returns any error that occurs.
func (r *InviteRepository) TestConnection() error {
	return r.mdb.TestConnection(r.dbName)
}

// Insert inserts the provided Invite into the database.
//
// inv is the Invite to insert into the database.
func (r *InviteRepository) Insert(inv *types.Invite) error {
	if err := r.mdb.InsertOne(r.dbName, r.collName, inv); err != nil {
		return fmt.Errorf("failed to insert invite: %v", err)
	}

	return nil
}

// Get returns the Invite for the provided email address and token (or nil if no such Invite exists).
//
// email is the email address of the Invite to look up.
//
// token is the token string that should match with the email on the Invite.
func (r *InviteRepository) Get(email string, token string) (*types.Invite, error) {
	// Define the query
	query := bson.M{"email": email}

	// Load Invite from the database
	var invs []types.Invite
	if err := r.mdb.GetAll(r.dbName, r.collName, query, &invs); err != nil {
		return nil, fmt.Errorf("failed to load invite")
	}

	// Verify that only one Invite was loaded
	if len(invs) == 0 {
		return nil, nil
	} else if len(invs) != 1 {
		return nil, fmt.Errorf("multiple invites exists for email=%s token=%s", email, token)
	}

	return &invs[0], nil
}

// Accept updates the Accepted property of the Invite keyed by the provided ID to be true
//
// id is the unique identifier of the Invite being accepted.
func (r *InviteRepository) Accept(id types.InviteID) error {
	// Define the filter query and update operation
	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"accepted": true,
		},
	}

	// Perform the update
	if err := r.mdb.UpdateOne(r.dbName, r.collName, filter, update); err != nil {
		return fmt.Errorf("failed to accept invite (id=%s): %v", id, err)
	}

	return nil
}
