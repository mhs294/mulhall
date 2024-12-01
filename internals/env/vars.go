package env

import (
	"fmt"
	"os"
)

var MongoDBConnStr string

func LoadVars() error {
	var err error
	MongoDBConnStr, err = loadVar("MULHALL_DB_CONN_STR")
	if err != nil {
		return err
	}

	return nil
}

func loadVar(name string) (string, error) {
	value := os.Getenv(name)
	if len(value) == 0 {
		return "", fmt.Errorf("missing required environment variable %s", name)
	}

	return value, nil
}
