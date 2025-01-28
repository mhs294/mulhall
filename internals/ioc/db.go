package ioc

import (
	"github.com/mhs294/mulhall/internals/db"
	"github.com/mhs294/mulhall/internals/env"
)

var mongoDB *db.MongoDB

func MongoDB() *db.MongoDB {
	if mongoDB == nil {
		logger = Logger()
		mongoDB = db.NewMongoDB(env.MongoDBConnStr, env.Timeout, logger)
	}

	return mongoDB
}
