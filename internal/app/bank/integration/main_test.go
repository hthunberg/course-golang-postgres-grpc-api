//go:build integration

package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hthunberg/course-golang-postgres-grpc-api/dbtest"
	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/app/bank"
)

var testee bank.Bank

func TestMain(m *testing.M) {
	os.Setenv("MIGRATION_URL", fmt.Sprintf("file:.%s", "/../../../../build/db/migrations"))

	ctx := context.Background()

	// Set up a postgres DB
	testDBRequest := dbtest.TestDatabaseContainerRequest()
	testDB, err := dbtest.SetupTestDatabase(ctx, testDBRequest)
	if err != nil {
		log.Fatal("failed to setup postgres db", err)
	}
	defer testDB.TearDown()
	testee = bank.NewBank(testDB.DbInstance)

	os.Exit(m.Run())
}
