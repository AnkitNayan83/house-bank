package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

const dbSource = "postgresql://root:postgres@localhost:5432/house_bank?sslmode=disable"

var testQueries *Queries
var testDb *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error

	testDb, err = pgxpool.New(context.Background(), dbSource)

	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}
	defer testDb.Close()

	testQueries = New(testDb)

	os.Exit(m.Run())
}
