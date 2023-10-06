//go:build integration

package integration

import (
	"fmt"
	"os"
	"testing"

	"github.com/hthunberg/course-golang-postgres-grpc-api/dbtest"
	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/app/bank"
)

var testee bank.Bank

func TestMain(m *testing.M) {
	os.Setenv("MIGRATION_URL", fmt.Sprintf("file:.%s", "/../../../../build/db/migrations"))

	testDB := dbtest.SetupTestDatabase()
	defer testDB.TearDown()
	testee = bank.NewBank(testDB.DbInstance)

	os.Exit(m.Run())
}
