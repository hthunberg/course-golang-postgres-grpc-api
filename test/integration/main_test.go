//go:build integration

package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/hthunberg/course-golang-postgres-grpc-api/dbtest"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testDbInstance  *pgxpool.Pool
	testBankBaseURL string
)

func TestMain(m *testing.M) {
	os.Setenv("MIGRATION_URL", fmt.Sprintf("file:.%s", "/../../build/db/migrations"))

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

	// Set up a postgres DB
	testDBRequest := dbtest.TestDatabaseContainerRequest()
	network.ApplyNetworkAlias(&testDBRequest, testDBAlias)
	testDB, err := dbtest.SetupTestDatabase(ctx, testDBRequest)
	if err != nil {
		log.Fatal("failed to setup postgres db", err)
	}
	defer testDB.TearDown()
	testDbInstance = testDB.DbInstance

	// Set up a test bank
	dbAddr := fmt.Sprintf("%s:%s", testDBAlias, dbtest.DBPort)
	testBankRequest := TestBankContainerRequest(dbAddr)
	network.ApplyNetworkAlias(&testBankRequest, "testbank")
	testBank, err := setupTestBank(ctx, testBankRequest)
	if err != nil {
		log.Fatal("failed to setup test bank", err)
	}
	defer testBank.TearDown()

	testBankBaseURL = testBank.URI

	time.Sleep(time.Second)

	os.Exit(m.Run())
}
