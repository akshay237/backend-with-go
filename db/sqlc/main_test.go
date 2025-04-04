package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/akshay237/backend-with-go/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("db connection failed", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
