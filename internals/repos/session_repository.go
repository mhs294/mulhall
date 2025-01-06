package repos

import (
	"github.com/mhs294/mulhall/internals/db"
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

// TODO - start here
