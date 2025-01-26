package repos

import (
	"fmt"

	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
)

// UserRepository is a mechanism for managing Users and their acccount details on the site.
type UserRepository struct {
	mdb      *db.MongoDB
	dbName   string
	collName string
}

// NewUserRepository creates a new UserRepository instance and returns a pointer to it.
//
// mdb is the MongoDB instance used by the UserRepository.
func NewUserRepository(mdb *db.MongoDB) *UserRepository {
	return &UserRepository{mdb: mdb, dbName: "mulhall", collName: "users"}
}

// TestConnection tests the connection to the MongoDB instance and returns any error that occurs.
func (r *UserRepository) TestConnection() error {
	return r.mdb.TestConnection(r.dbName)
}

// InsertUser inserts the provided User into the database.
//
// u is the User to insert into the database.
func (r *UserRepository) InsertUser(u *types.User) error {
	if err := r.mdb.InsertOne(r.dbName, r.collName, u); err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	return nil
}

// GetUser returns the User for the provided email address (or nil if no such User exists).
//
// email is the email address of the User to look up.
func (r *UserRepository) GetUser(email string) (*types.User, error) {
	// Load User from the database
	var users []types.User
	if err := r.mdb.GetAll(r.dbName, r.collName, bson.D{{Key: "email", Value: email}}, &users); err != nil {
		return nil, fmt.Errorf("failed to look up user: %v", err)
	}

	// Verify that only one User was loaded
	if len(users) == 0 {
		return nil, &types.UserNotFoundError{Email: email}
	} else if len(users) != 1 {
		return nil, fmt.Errorf("multiple users exists for email=%s", email)
	}

	return &users[0], nil
}
