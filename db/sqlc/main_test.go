package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:password@localhost:5432/simple_bank?sslmode=disable&timezone=UTC&connect_timeout=5"
)

//Setup db connection profile
var testQueries *Queries

// Setup reusable db object
var testDB *sql.DB

// (TestMain) This is the main entry point for all tests within this package (db)
func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db :=", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
