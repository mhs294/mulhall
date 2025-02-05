package repos

import (
	"fmt"

	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
)

// UserRepository manages User records in the database.
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

// Insert inserts the provided User into the database.
//
// u is the User to insert into the database.
func (r *UserRepository) Insert(u *types.User) error {
	if err := r.mdb.InsertOne(r.dbName, r.collName, u); err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	return nil
}

// GetByEmail returns the active User for the provided email address.
// Returns UserNotFoundError if no such User exists.
// Returns UserInactiveError if a User exists for the email address but is inactive.
//
// email is the email address of the User to look up.
func (r *UserRepository) GetByEmail(email string) (*types.User, error) {
	// Define the query
	query := bson.M{"email": email}

	// Load User from the database
	var users []types.User
	if err := r.mdb.GetAll(r.dbName, r.collName, query, &users); err != nil {
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

// GetByID returns the active User with the provided ID.
// Returns UserNotFoundError if no such User exists or the User has been deactivated.
//
// id is the unique identifier of the User to look up.
func (r *UserRepository) GetByID(id types.UserID) (*types.User, error) {
	// Define the query
	query := bson.M{"id": id}

	// Load User from the database
	var u types.User
	if err := r.mdb.GetOne(r.dbName, r.collName, query, &u); err != nil {
		return nil, fmt.Errorf("failed to look up user: %v", err)
	}

	// Verify that the User exists and is active
	if u.ID != id || !u.Active {
		return nil, &types.UserNotFoundError{ID: id}
	}

	return &u, nil
}
