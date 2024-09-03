package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/PaBah/GophKeeper/db"
	"github.com/PaBah/GophKeeper/internal/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
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

// CreateUser - create new user
func (ds *DBStorage) CreateUser(ctx context.Context, user models.User) (createdUser models.User, err error) {
	_, DBerr := ds.db.ExecContext(ctx,
		`INSERT INTO users(email, password) VALUES ($1, $2)`, user.Email, user.Password)

	var pgErr *pgconn.PgError
	if errors.As(DBerr, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		err = ErrAlreadyExists
		return
	}

	return ds.AuthorizeUser(ctx, user.Email)
}

func (ds *DBStorage) AuthorizeUser(ctx context.Context, email string) (user models.User, err error) {
	row := ds.db.QueryRowContext(ctx, `SELECT id, password FROM users WHERE email=$1`, email)
	var id, password string
	err = row.Scan(&id, &password)

	if err != nil {
		return
	}

	user = models.User{ID: id, Email: email, Password: password}
	return
}

// Close - close connection to Data Base
func (ds *DBStorage) Close() error {
	return ds.db.Close()
}

// NewDBStorage - create instance of DBStorage
func NewDBStorage(ctx context.Context, databaseDSN string) (DBStorage, error) {
	store := DBStorage{}
	err := store.initialize(ctx, databaseDSN)
	return store, err
}
