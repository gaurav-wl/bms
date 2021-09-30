package db

import (
	"fmt"
	"github.com/gauravcoco/bms/providers"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"log"
)

type Migrations struct {
	DB *sqlx.DB
}

func NewMigrationProvider(db *sqlx.DB) providers.MigrationProvider {
	return &Migrations{
		DB: db,
	}
}

func (m Migrations) Up() {
	driver, err := postgres.WithInstance(m.DB.DB, &postgres.Config{})
	migration, err := migrate.NewWithDatabaseInstance("file://db/migrations/", "dcnom9cktasj0k", driver)
	if err != nil {
		logrus.Fatalf("Unable to fetch Migrations %v", err)
	}
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Println("Unable to run Migrations")
		log.Fatal(err)
	}
	fmt.Println("Migration Up Successful")
}
