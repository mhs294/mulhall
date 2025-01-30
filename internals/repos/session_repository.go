package repos

import (
	"fmt"
	"time"

	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/types"
	"go.mongodb.org/mongo-driver/bson"
)

// SessionRepository is a mechanism for managing Sessions of logged in Users on the site.
type SessionRepository struct {
	mdb      *db.MongoDB
	dbName   string
	collName string
}

// NewSessionRepository creates a new SessionRepository instance and returns a pointer to it.
//
// mdb is the MongoDB instance used by the SessionRepository.
func NewSessionRepository(mdb *db.MongoDB) *SessionRepository {
	return &SessionRepository{mdb: mdb, dbName: "mulhall", collName: "sessions"}
}

// TestConnection tests the connection to the MongoDB instance and returns any error that occurs.
func (r *SessionRepository) TestConnection() error {
	return r.mdb.TestConnection(r.dbName)
}

// Insert inserts the provided Session into the database.
//
// s is the Session to insert into the database.
func (r *SessionRepository) Insert(s *types.Session) error {
	if err := r.mdb.InsertOne(r.dbName, r.collName, s); err != nil {
		return fmt.Errorf("failed to insert session: %v", err)
	}

	return nil
}

// GetByID returns the Session for the provided ID (or nil if no such Session exists).
//
// id is the unique identifier of the Session to look up.
func (r *SessionRepository) GetByID(id types.SessionID) (*types.Session, error) {
	// Define the query
	query := bson.M{"id": id}

	// Load the Session from the database
	var sess types.Session
	if err := r.mdb.GetOne(r.dbName, r.collName, query, &sess); err != nil {
		return nil, fmt.Errorf("failed to look up session: %v", err)
	}

	// Verify that the Session exists and is active (i.e. - has not expired)
	if sess == (types.Session{}) {
		return nil, &types.SessionNotFoundError{ID: id}
	} else if sess.Expiration.Before(time.Now().UTC()) {
		return nil, &types.SessionExpiredError{}
	}

	return &sess, nil
}
