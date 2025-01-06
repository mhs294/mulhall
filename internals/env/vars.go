package env

import (
	"fmt"
	"os"
	"time"
)

var MongoDBConnStr string
var Timeout time.Duration
var InviteExpiration time.Duration
var SessionExpiration time.Duration

func LoadVars() error {
	var err error
	MongoDBConnStr, err = loadVar("MULHALL_DB_CONN_STR")
	if err != nil {
		return err
	}

	// TODO - make these configurable
	Timeout = time.Second * 10
	InviteExpiration = time.Hour * 24 * 7
	SessionExpiration = time.Hour * 24 * 7

	return nil
}

func loadVar(name string) (string, error) {
	value := os.Getenv(name)
	if len(value) == 0 {
		return "", fmt.Errorf("missing required environment variable %s", name)
	}

	return value, nil
}
