//go:build integration

package integration

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/hthunberg/course-golang-postgres-grpc-api/dbtest"
	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/app/bank"
	"go.uber.org/zap"
)

var testee bank.Bank

func TestMain(m *testing.M) {
	ctx := context.Background()

	absoluteMigrationsPath, err := filepath.Abs("../../../../build/db/migrations")
	if err != nil {
		log.Fatal("failed to calculate absolute path to db migrations", zap.Error(err))
	}

	// Set up a postgres DB
	testDBRequest := dbtest.TestDatabaseContainerRequest()
	testDB, err := dbtest.SetupTestDatabase(ctx, testDBRequest, absoluteMigrationsPath)
	if err != nil {
		log.Fatal("failed to setup postgres db", zap.Error(err))
	}
	defer testDB.TearDown()
	testee = bank.NewBank(testDB.DbInstance)

	os.Exit(m.Run())
}
