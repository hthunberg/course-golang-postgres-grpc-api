//go:build integration

package integration

import (
	"fmt"
	"os"
	"testing"

	"github.com/hthunberg/course-golang-postgres-grpc-api/dbtest"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testDbInstance *pgxpool.Pool

func TestMain(m *testing.M) {
	os.Setenv("MIGRATION_URL", fmt.Sprintf("file:.%s", "/../../build/db/migrations"))

	testDB := dbtest.SetupTestDatabase()
	defer testDB.TearDown()
	testDbInstance = testDB.DbInstance
	os.Exit(m.Run())
}
