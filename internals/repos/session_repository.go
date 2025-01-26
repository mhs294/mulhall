package repos

import (
	"fmt"

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

// InsertSession inserts the provided Session into the database.
//
// s is the Session to insert into the database.
func (r *SessionRepository) InsertSession(s *types.Session) error {
	if err := r.mdb.InsertOne(r.dbName, r.collName, s); err != nil {
		return fmt.Errorf("failed to insert session: %v", err)
	}

	return nil
}

// GetSession returns the Session for the provided ID (or nil if no such Session exists).
//
// id is the unique identifier of the Session to look up.
func (r *SessionRepository) GetSession(id types.SessionID) (*types.Session, error) {
	var sess *types.Session
	if err := r.mdb.GetOne(r.dbName, r.collName, bson.D{{Key: "id", Value: id}}, &sess); err != nil {
		return nil, fmt.Errorf("failed to look up session: %v", err)
	}
	if sess == nil {
		return nil, &types.SessionNotFoundError{ID: id}
	}

	return sess, nil
}
