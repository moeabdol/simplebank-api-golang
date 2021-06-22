package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testStore *Store

const (
	dbDriver = "postgres"
	dbSource = "postgres://postgres:postgres@localhost:5432/simplebank_test_db?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database")
	}

	testStore = NewStore(conn)

	os.Exit(m.Run())
}
