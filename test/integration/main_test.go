//go:build integration

package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hthunberg/course-golang-postgres-grpc-api/dbtest"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	testDbInstance  *pgxpool.Pool
	testBankBaseURL string
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Docker provides the ability for us to create custom networks and place containers on one or more networks.
	// The communication can then occur between networked containers without the need of exposing ports through the host.
	// In this particular setup TestBank can access the TestDB using network alias "testdb" and its internal port 5432.
	network, err := CreateNetwork(ctx)
	if err != nil {
		log.Fatal("failed to setup test network", err)
	}
	defer network.TearDown(ctx)

	// Alias for postgres db when running inside a custom network,
	testDBAlias := "testdb"

	absoluteMigrationsPath, err := filepath.Abs("../../build/db/migrations")
	if err != nil {
		log.Fatal("failed to calculate absolute path to db migrations", zap.Error(err))
	}

	// Set up a postgres DB
	testDBRequest := dbtest.TestDatabaseContainerRequest()
	network.ApplyNetworkAlias(&testDBRequest, testDBAlias)
	testDB, err := dbtest.SetupTestDatabase(ctx, testDBRequest, absoluteMigrationsPath)
	if err != nil {
		log.Fatal("failed to setup postgres db", zap.Error(err))
	}
	defer testDB.TearDown()
	testDbInstance = testDB.DbInstance

	// Set up a test bank
	dbAddr := fmt.Sprintf("%s:%s", testDBAlias, dbtest.DBPort)
	testBankRequest := TestBankContainerRequest(dbAddr, absoluteMigrationsPath)
	network.ApplyNetworkAlias(&testBankRequest, "testbank")
	testBank, err := setupTestBank(ctx, testBankRequest)
	if err != nil {
		log.Fatal("failed to setup test bank", zap.Error(err))
	}
	defer testBank.TearDown()

	testBankBaseURL = testBank.URI

	time.Sleep(time.Second)

	os.Exit(m.Run())
}
