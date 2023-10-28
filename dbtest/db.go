package dbtest

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // used by golang-migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // used by golang-migrate
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // used by golang-migrate
)

const (
	DBName = "test_db"
	DBUser = "test_user"
	DBPass = "test_password"
	DBPort = "5432"
	port   = DBPort + "/tcp"
)

// TestDatabase represents
// - connection pool, a pool of connections ready to use
// - db address (host:port) to the running db
// - handle to running test container
type TestDatabase struct {
	DbInstance *pgxpool.Pool
	DBPort     string
	DBHost     string
	Container  testcontainers.Container
}

func SetupTestDatabase(ctx context.Context, testDatabaseContainerRequest testcontainers.GenericContainerRequest, absoluteMigrationsPath string) (*TestDatabase, error) {
	// setup db container
	container, dbInstance, err := createContainer(ctx, testDatabaseContainerRequest)
	if err != nil {
		return nil, fmt.Errorf("setup test db:create container: %v", err)
	}

	dbPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return nil, fmt.Errorf("setup test bank:container port: %v", err)
	}

	dbHost, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("setup test db:container host: %v", err)
	}

	dbAddr := fmt.Sprintf("%s:%s", dbHost, dbPort.Port())

	// migrate db schema
	err = migrateDb(dbAddr, absoluteMigrationsPath)
	if err != nil {
		return nil, fmt.Errorf("setup test bank:migrate db: %v", err)
	}

	return &TestDatabase{
		Container:  container,
		DbInstance: dbInstance,
		DBPort:     dbPort.Port(),
		DBHost:     dbHost,
	}, nil
}

// TearDown tears down the running database container
func (tdb *TestDatabase) TearDown() {
	tdb.DbInstance.Close()
	// remove test container
	_ = tdb.Container.Terminate(context.Background())
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

	dbAddr := fmt.Sprintf("%s:%s", tdb.DBHost, tdb.DBPort)

	log.Println("database truncated: ", dbAddr)
	return nil
}

func TestDatabaseContainerRequest() testcontainers.GenericContainerRequest {
	env := map[string]string{
		"POSTGRES_PASSWORD": DBPass,
		"POSTGRES_USER":     DBUser,
		"POSTGRES_DB":       DBName,
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:bullseye",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor: wait.ForAll(
				wait.ForLog("database system is ready to accept connections"),
				wait.ForExposedPort().WithStartupTimeout(60*time.Second),
				wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(10*time.Second),
			),
		},
		Started: true,
	}

	return req
}

func createContainer(ctx context.Context, req testcontainers.GenericContainerRequest) (testcontainers.Container, *pgxpool.Pool, error) {
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, fmt.Errorf("create container:failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return container, nil, fmt.Errorf("create container:failed to get container external port: %v", err)
	}

	h, err := container.Host(ctx)
	if err != nil {
		return container, nil, fmt.Errorf("create container:failed to get container host: %v", err)
	}

	time.Sleep(time.Second)

	dbAddr := fmt.Sprintf("%s:%s", h, p.Port())

	log.Println("postgres container ready and running at: ", dbAddr)

	db, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DBUser, DBPass, dbAddr, DBName))
	if err != nil {
		return container, db, fmt.Errorf("create container:failed to establish database connection: %v", err)
	}

	return container, db, nil
}

func migrateDb(dbAddr, absoluteMigrationsPath string) error {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DBUser, DBPass, dbAddr, DBName)

	// URL file:./<path> needed fir golang-migrate
	migrationsURL := fmt.Sprintf("file:%s", absoluteMigrationsPath)

	if len(migrationsURL) < 1 {
		return fmt.Errorf("migrate db:missing env migration_url: %s", migrationsURL)
	}

	log.Printf("migrate db:running db migrations using db %s migrations %s", databaseURL, migrationsURL)

	migration, err := migrate.New(migrationsURL, databaseURL)
	if err != nil {
		return err
	}
	defer migration.Close()

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate db:migrate up: %v", err)
	}

	log.Println("migration done")

	return nil
}
