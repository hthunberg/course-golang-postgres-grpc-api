package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/hthunberg/course-golang-postgres-grpc-api/cmd"
	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/app/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	_ "github.com/golang-migrate/migrate/v4/database/postgres" // golang-migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // golang-migrate
	_ "github.com/jackc/pgx/v5/stdlib"                         // golang-migrate
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := util.LoadConfig("./config")
	if err != nil {
		// We need to use the default logger since zap is not ready yet
		log.Fatalln("initializing:", err)
	}

	logger, err := cmd.NewLogger(cfg.LogLevel)
	if err != nil {
		// We need to use the default logger since zap is not ready yet
		log.Fatalln("initializing:", err)
	}

	// Flush any buffered logs when exiting main
	defer logger.Sync()

	logger.Info(
		"initializing: starting application",
		zap.String("build_version", "0.0.1"),
		zap.String("environment", cfg.Environment),
	)

	connPool, err := pgxpool.New(ctx, cfg.DBSource)
	if err != nil {
		logger.Fatal("initializing: connect to db", zap.Error(err))
	}

	if err := connPool.Ping(ctx); err != nil {
		logger.Fatal("initializing: ping db", zap.Error(err))
	}

	runDBMigration(cfg.MigrationURL, cfg.DBSource, logger)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-termChan:
		logger.Info("context cancel: closing application")
		cancel()
	case <-ctx.Done():
		logger.Info("context done")
	}
}

func runDBMigration(migrationURL string, dbSource string, logger *zap.Logger) {
	logger.Info("initializing: migrate db", zap.String("migrations", migrationURL))

	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		logger.Fatal(
			"initializing: migrate db: new instance",
			zap.String("migrations", migrationURL),
			zap.Error(err),
		)
	}

	defer migration.Close()

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal("initializing: migrate db: migrate up", zap.Error(err))
	}

	logger.Info("initializing: migrate db: finished")
}
