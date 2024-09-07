package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/PaBah/GophKeeper/db"
	"github.com/PaBah/GophKeeper/internal/config"
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

// CreateCredentials - create new credentials record
func (ds *DBStorage) CreateCredentials(ctx context.Context, credentials models.Credentials) (createdCredentials models.Credentials, err error) {
	createdCredentials = credentials
	_, DBerr := ds.db.ExecContext(ctx,
		`INSERT INTO credentials(service_name, identity, password, user_id) VALUES ($1, $2, $3, $4)`,
		credentials.ServiceName, credentials.Identity, credentials.Password, ctx.Value(config.USERIDCONTEXTKEY).(string))

	var pgErr *pgconn.PgError
	if errors.As(DBerr, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		err = ErrAlreadyExists
		return
	}

	row := ds.db.QueryRowContext(ctx, `SELECT id, uploaded_at FROM credentials WHERE service_name=$1 and user_id=$2`,
		credentials.ServiceName, ctx.Value(config.USERIDCONTEXTKEY).(string))

	_ = row.Scan(&createdCredentials.ID, &createdCredentials.UploadedAt)

	return
}

// GetCredentials - return list of users Credentials
func (ds *DBStorage) GetCredentials(ctx context.Context) (credentials []models.Credentials, err error) {
	var rows *sql.Rows
	rows, err = ds.db.QueryContext(ctx,
		`SELECT id, service_name, identity, password, uploaded_at FROM credentials WHERE user_id=$1`,
		ctx.Value(config.USERIDCONTEXTKEY).(string))
	if err != nil {
		return
	}
	err = rows.Err()
	defer rows.Close()

	credentials = make([]models.Credentials, 0)
	for rows.Next() {
		var credentialSet models.Credentials
		err = rows.Scan(&credentialSet.ID, &credentialSet.ServiceName, &credentialSet.Identity, &credentialSet.Password, &credentialSet.UploadedAt)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, credentialSet)
	}
	return
}

// UpdateCredentials - update Credentials model
func (ds *DBStorage) UpdateCredentials(ctx context.Context, credentials models.Credentials) (updatedCredentials models.Credentials, err error) {
	_, err = ds.db.ExecContext(ctx,
		`UPDATE credentials SET service_name=$1, identity=$2, password=$3 WHERE user_id=$4 and id=$5`,
		credentials.ServiceName, credentials.Identity, credentials.Password, ctx.Value(config.USERIDCONTEXTKEY).(string), credentials.ID)
	if err == nil {
		updatedCredentials = credentials
	}

	row := ds.db.QueryRowContext(ctx, `SELECT uploaded_at FROM credentials WHERE id=$1 and user_id=$2`,
		credentials.ID, ctx.Value(config.USERIDCONTEXTKEY).(string))

	_ = row.Scan(&updatedCredentials.UploadedAt)

	return
}

// DeleteCredentials - update Credentials model
func (ds *DBStorage) DeleteCredentials(ctx context.Context, credentialsID string) (err error) {
	_, err = ds.db.ExecContext(ctx,
		`DELETE FROM credentials WHERE user_id=$1 and id=$2`, ctx.Value(config.USERIDCONTEXTKEY).(string), credentialsID)
	return
}

// CreateCard - create new Card record
func (ds *DBStorage) CreateCard(ctx context.Context, card models.Card) (createdCard models.Card, err error) {
	createdCard = card
	_, DBerr := ds.db.ExecContext(ctx,
		`INSERT INTO cards(number, expiration_date, holder_name, cvv, user_id) VALUES ($1, $2, $3, $4, $5)`,
		card.Number, card.ExpirationDate, card.HolderName, card.CVV,
		ctx.Value(config.USERIDCONTEXTKEY).(string))

	var pgErr *pgconn.PgError
	if errors.As(DBerr, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		err = ErrAlreadyExists
		return
	}

	row := ds.db.QueryRowContext(ctx, `SELECT id, uploaded_at FROM cards WHERE number=$1 and user_id=$2`,
		card.Number, ctx.Value(config.USERIDCONTEXTKEY).(string))

	_ = row.Scan(&createdCard.ID, &createdCard.UploadedAt)

	return
}

// GetCards - return list of users Cards
func (ds *DBStorage) GetCards(ctx context.Context) (cards []models.Card, err error) {
	var rows *sql.Rows
	rows, err = ds.db.QueryContext(ctx,
		`SELECT id, number, expiration_date, holder_name, cvv, uploaded_at FROM cards WHERE user_id=$1`,
		ctx.Value(config.USERIDCONTEXTKEY).(string))
	if err != nil {
		return
	}
	err = rows.Err()
	defer rows.Close()

	cards = make([]models.Card, 0)
	for rows.Next() {
		var card models.Card
		err = rows.Scan(&card.ID, &card.Number, &card.ExpirationDate, &card.HolderName, &card.CVV, &card.UploadedAt)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return
}

// UpdateCard - update Card model
func (ds *DBStorage) UpdateCard(ctx context.Context, card models.Card) (updatedCard models.Card, err error) {

	_, err = ds.db.ExecContext(ctx,
		`UPDATE cards SET number=$1, expiration_date=$2, holder_name=$3, cvv=$4 WHERE user_id=$5 and id=$6`,
		card.Number, card.ExpirationDate, card.HolderName, card.CVV, ctx.Value(config.USERIDCONTEXTKEY).(string), card.ID)
	if err == nil {
		updatedCard = card
	}

	row := ds.db.QueryRowContext(ctx, `SELECT uploaded_at FROM cards WHERE id=$1 and user_id=$2`,
		card.ID, ctx.Value(config.USERIDCONTEXTKEY).(string))

	_ = row.Scan(&updatedCard.UploadedAt)

	return
}

// DeleteCard - update Card model
func (ds *DBStorage) DeleteCard(ctx context.Context, cardID string) (err error) {
	_, err = ds.db.ExecContext(ctx,
		`DELETE FROM cards WHERE user_id=$1 and id=$2`, ctx.Value(config.USERIDCONTEXTKEY).(string), cardID)
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
