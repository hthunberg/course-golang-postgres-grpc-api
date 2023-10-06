package dbtest

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // used by golang-migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // used by golang-migrate
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // used by golang-migrate
)

const (
	DbName = "test_db"
	DbUser = "test_user"
	DbPass = "test_password"
)

// TestDatabase represents
// - connection pool, a pool of connections ready to use
// - db address (host:port) to the running db
// - handle to running test container
type TestDatabase struct {
	DbInstance *pgxpool.Pool
	DbAddress  string
	container  testcontainers.Container
}

func SetupTestDatabase() *TestDatabase {
	// setup db container
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	container, dbInstance, dbAddr, err := createContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}

	// migrate db schema
	err = migrateDb(dbAddr)
	if err != nil {
		log.Fatal("failed to perform db migration", err)
	}
	cancel()

	return &TestDatabase{
		container:  container,
		DbInstance: dbInstance,
		DbAddress:  dbAddr,
	}
}

// TearDown tears down the running database container
func (tdb *TestDatabase) TearDown() {
	tdb.DbInstance.Close()
	// remove test container
	_ = tdb.container.Terminate(context.Background())
}

func (tdb *TestDatabase) Truncate() error {
	query := []string{
		"TRUNCATE accounts CASCADE",
		"TRUNCATE entries CASCADE",
		"TRUNCATE transfers CASCADE",
		"TRUNCATE users CASCADE",
	}
	for _, q := range query {
		_, err := tdb.DbInstance.Exec(context.Background(), q)
		if err != nil {
			return fmt.Errorf("failed to truncate db: %v", err)
		}
	}

	log.Println("database truncated: ", tdb.DbAddress)
	return nil
}

func createContainer(ctx context.Context) (testcontainers.Container, *pgxpool.Pool, string, error) {
	env := map[string]string{
		"POSTGRES_PASSWORD": DbPass,
		"POSTGRES_USER":     DbUser,
		"POSTGRES_DB":       DbName,
	}
	port := "5432/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:bullseye",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container external port: %v", err)
	}

	h, err := container.Host(ctx)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container host: %v", err)
	}

	time.Sleep(time.Second)

	dbAddr := fmt.Sprintf("%s:%s", h, p.Port())

	log.Println("postgres container ready and running at: ", dbAddr)

	db, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName))
	if err != nil {
		return container, db, dbAddr, fmt.Errorf("failed to establish database connection: %v", err)
	}

	return container, db, dbAddr, nil
}

func migrateDb(dbAddr string) error {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName)

	// file:./<path> to be relative to working directory
	migrationsURL := os.Getenv("MIGRATION_URL")

	if len(migrationsURL) < 1 {
		return fmt.Errorf("missing env migration_url: %s", migrationsURL)
	}

	migration, err := migrate.New(migrationsURL, databaseURL)
	if err != nil {
		return err
	}
	defer migration.Close()

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("migration done")

	return nil
}
