package storage

import (
	"context"
	"database/sql"

	"github.com/PaBah/GophKeeper/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// DBStorage - model of Repository storage on top of Data Base
type DBStorage struct {
	db *sql.DB
}

func (ds *DBStorage) initialize(ctx context.Context, databaseDSN string) (err error) {

	ds.db, err = sql.Open("pgx", databaseDSN)
	if err != nil {
		return
	}

	driver, err := iofs.New(db.MigrationsFS, "migrations")
	if err != nil {
		return err
	}

	d, err := postgres.WithInstance(ds.db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", driver, "psql_db", d)
	if err != nil {
		return err
	}

	_ = m.Up()
	return
}

// NewDBStorage - create instance of DBStorage
func NewDBStorage(ctx context.Context, databaseDSN string) (DBStorage, error) {
	store := DBStorage{}
	err := store.initialize(ctx, databaseDSN)
	return store, err
}
