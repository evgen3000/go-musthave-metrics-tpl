package postgres

import (
	"context"
	"log"

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
		log.Fatalf("Unable to connect to database: %v", err)
	}

	q := `create table gauge
			(id varchar(256) primary key,
			value double precision);
		create table counter (
		    id varchar(256) primary key,
		    value integer)`
	_, errExec := pool.Exec(context.Background(), q)
	if errExec != nil {
		log.Fatalf("Unable to create table: %v", err)
	}
	pingErr := pool.Ping(context.Background())
	if pingErr != nil {
		log.Fatalf("Unable to ping database: %v", pingErr)
	}
	log.Println("Successfully connected to database")
	return pool
}
