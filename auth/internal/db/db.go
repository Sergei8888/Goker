package db

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"os"
)

type dbConfig struct {
	User     string
	Password string
	Dbname   string
	Host     string
}

func MustInitDb() *sqlx.DB {
	config := mustLoadDbConfig()
	db := mustConnectDb(config)
	mustMigrateDb(db, config)

	return db
}

func mustLoadDbConfig() dbConfig {
	user, found := os.LookupEnv("PG_USER")
	if found == false {
		panic("PG_USER env variable not set")
	}

	password, found := os.LookupEnv("PG_PASSWORD")
	if found == false {
		panic("PG_PASSWORD env variable not set")
	}

	dbname, found := os.LookupEnv("PG_DBNAME")
	if found == false {
		panic("PG_DBNAME env variable not set")
	}

	host, found := os.LookupEnv("PG_HOST")
	if found == false {
		panic("PG_URL env variable not set")
	}

	return dbConfig{user, password, dbname, host}
}

func mustConnectDb(config dbConfig) *sqlx.DB {

	connectionString := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s sslmode=disable",
		config.User,
		config.Password,
		config.Dbname,
		config.Host)

	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		panic(err)
	}

	return db
}

func mustMigrateDb(db *sqlx.DB, config dbConfig) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations",
		config.Dbname, driver)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
}
