package repos

import (
	"fmt"

	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
)

// InviteRepository is a mechanism for managing invitations for new Users to join the site.
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
func (ir *InviteRepository) TestConnection() error {
	return ir.mdb.TestConnection(ir.dbName)
}

// InsertInvite inserts the provided Invite into the database.
//
// inv is the Invite to insert into the database.
func (ir *InviteRepository) InsertInvite(inv *types.Invite) error {
	// Insert the BSON map into the invites collection as a document
	if err := ir.mdb.InsertOne(ir.dbName, ir.collName, inv); err != nil {
		return fmt.Errorf("failed to insert invite: %v", err)
	}

	return nil
}

// GetInvite returns the Invite for the provided email address and token (or nil if no such Invite exists).
//
// email is the email address to look up the Invite for.
//
// token is the token string that should match with the email on the Invite.
func (ir *InviteRepository) GetInvite(email string, token string) (*types.Invite, error) {
	// Load invite from the database
	var invites []types.Invite
	if err := ir.mdb.GetAll(ir.dbName, ir.collName, bson.D{{Key: "email", Value: email}}, &invites); err != nil {
		return nil, fmt.Errorf("failed to load invite")
	}

	// Verify that only one Invite was loaded
	if len(invites) == 0 {
		return nil, nil
	} else if len(invites) != 1 {
		return nil, fmt.Errorf("multiple invites exists for email=%s token=%s", email, token)
	}

	return &invites[0], nil
}

// AcceptInvite updates the Accepted property of the Invite keyed by the provided ID to be true
//
// id is the unique identifier of the Invite being accepted.
func (ir *InviteRepository) AcceptInvite(id types.InviteID) error {
	// Define the update operation and filter
	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"accepted": true,
		},
	}

	// Perform the update
	if err := ir.mdb.UpdateOne(ir.dbName, ir.collName, filter, update); err != nil {
		return fmt.Errorf("failed to accept invite (id=%s): %v", id, err)
	}

	return nil
}
