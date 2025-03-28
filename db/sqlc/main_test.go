package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/AnkitNayan83/houseBank/util"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	testDb, err = pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}
	defer testDb.Close()

	testQueries = New(testDb)

	os.Exit(m.Run())
}
