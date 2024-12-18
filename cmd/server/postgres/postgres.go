package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

//postgresql://admin:admin@localhost:5432/admin?schema=public

func Connect(dsn string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}

	pool, poolErr := pgxpool.NewWithConfig(context.Background(), config)
	if poolErr != nil {
		log.Fatalf("Unable to connect to database: %v", poolErr)
	}

	db, err := pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Unable to acquire a connection: %v", err)
	}
	defer db.Release()

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		log.Fatalf("Unable to create migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/server/postgres/migrations",
		"admin",
		driver,
	)
	if err != nil {
		log.Fatalf("Unable to create migrate instance: %v", err)
	}

	migrationErr := m.Up()
	if migrationErr != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migration failed: %v", err)
	}

	pingErr := pool.Ping(context.Background())
	if pingErr != nil {
		log.Fatalf("Unable to ping database: %v", pingErr)
	}
	log.Println("Successfully connected to database")

	return pool
}
