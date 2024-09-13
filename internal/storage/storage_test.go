package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PaBah/GophKeeper/internal/config"
	"github.com/PaBah/GophKeeper/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDBStorage(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		databaseDSN string
		wantErr     bool
	}{
		{
			name:        "Valid DSN",
			databaseDSN: "test",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st, err := NewDBStorage(ctx, tt.databaseDSN)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewDBStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				require.NotNil(t, st, "DBStorage should not be nil")
			}
		})
	}
}

func TestDBStorage_DeleteShortURLs(t *testing.T) {
	db, mock, _ := sqlmock.New()
	ds := &DBStorage{
		db: db,
	}
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM cards WHERE user_id=$1 and id=$2`)).
		WithArgs("test", "test").
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)
	err := ds.DeleteCard(context.WithValue(context.Background(), config.USERIDCONTEXTKEY, "test"), "test")
	assert.NoError(t, err, "successfully deleted card")
}

func TestDBStorage_UpdateCard(t *testing.T) {
	db, mock, _ := sqlmock.New()
	ds := &DBStorage{
		db: db,
	}
	ctx := context.WithValue(context.Background(), config.USERIDCONTEXTKEY, "test")
	card := models.NewCard("1234 5678 9012 3456", "12/24", "Test User", "123")
	card.ID = "1"
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE cards SET number=$1, expiration_date=$2, holder_name=$3, cvv=$4 WHERE user_id=$5 and id=$6`)).
		WithArgs(card.Number, card.ExpirationDate, card.HolderName, card.CVV, ctx.Value(config.USERIDCONTEXTKEY).(string), "1").
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT uploaded_at FROM cards WHERE id=$1 and user_id=$2`)).
		WithArgs(1, "test").
		WillReturnRows(sqlmock.NewRows([]string{"uploaded_at"}).
			AddRow(time.Now()))
	_, err := ds.UpdateCard(ctx, card)
	assert.NoError(t, err, "successfully updated card")
}

func TestDBStorage_UpdateCard_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	ds := &DBStorage{
		db: db,
	}
	ctx := context.WithValue(context.Background(), config.USERIDCONTEXTKEY, "test")
	card := models.NewCard("1234 5678 9012 3456", "12/24", "Test User", "123")
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE cards SET number=$1, expiration_date=$2, holder_name=$3, cvv=$4 WHERE user_id=$5 and id=$6`)).
		WithArgs(card.Number, card.ExpirationDate, card.HolderName, card.CVV, ctx.Value(config.USERIDCONTEXTKEY).(string), 1).
		WillReturnError(fmt.Errorf("an error"))
	_, err := ds.UpdateCard(ctx, card)
	assert.NotNil(t, err, "error should occur")
}

func TestDBStorage_UpdateCard_NoRowsAffected(t *testing.T) {
	db, mock, _ := sqlmock.New()
	ds := &DBStorage{
		db: db,
	}
	ctx := context.WithValue(context.Background(), config.USERIDCONTEXTKEY, "test")
	card := models.NewCard("1234 5678 9012 3456", "12/24", "Test User", "123")
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE cards SET number=$1, expiration_date=$2, holder_name=$3, cvv=$4 WHERE user_id=$5 and id=$6`)).
		WithArgs(card.Number, card.ExpirationDate, card.HolderName, card.CVV, ctx.Value(config.USERIDCONTEXTKEY).(string), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	_, err := ds.UpdateCard(ctx, card)
	assert.NotNil(t, err, "error should occur because no rows were affected")
}

func TestDBStorage_UpdateCard_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	ds := &DBStorage{
		db: db,
	}
	ctx := context.WithValue(context.Background(), config.USERIDCONTEXTKEY, "test")
	card := models.NewCard("1234 5678 9012 3456", "12/24", "Test User", "123")
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE cards SET number=$1, expiration_date=$2, holder_name=$3, cvv=$4 WHERE user_id=$5 and id=$6`)).
		WithArgs(card.Number, card.ExpirationDate, card.HolderName, card.CVV, ctx.Value(config.USERIDCONTEXTKEY).(string), 1).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT uploaded_at FROM cards WHERE id=$1 and user_id=$2`)).
		WithArgs(1, "test").
		WillReturnRows(sqlmock.NewRows([]string{"uploaded_at"}).
			AddRow("string-instead-of-time"))
	_, err := ds.UpdateCard(ctx, card)
	assert.NotNil(t, err, "should throw error on incorrect row scan")
}

func TestDBStorage_GetCards(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*DBStorage)
		cardCount int
		wantErr   bool
	}{
		{
			name:      "zero cards",
			setup:     func(ds *DBStorage) {},
			cardCount: 0,
			wantErr:   false,
		},
		{
			name:      "one card",
			setup:     func(ds *DBStorage) {},
			cardCount: 1,
			wantErr:   false,
		},
		{
			name:      "multiple cards",
			setup:     func(ds *DBStorage) {},
			cardCount: 5,
			wantErr:   false,
		},
		{
			name:      "error when getting cards",
			setup:     func(ds *DBStorage) {},
			cardCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			ds := &DBStorage{db: db}

			tt.setup(ds)
			if tt.wantErr {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, number, expiration_date, holder_name, cvv, uploaded_at FROM cards WHERE user_id=$1`)).
					WithArgs("test").
					WillReturnError(errors.New("some error"))
				return
			}

			rows := sqlmock.NewRows([]string{"id", "number", "expiration_date", "holder_name", "cvv", "uploaded_at"})
			for i := 0; i < tt.cardCount; i++ {
				rows.AddRow("test_id", "1234567812345678", "12/34", "Test User", "123", time.Now())
			}

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, number, expiration_date, holder_name, cvv, uploaded_at FROM cards WHERE user_id=$1`)).
				WithArgs("test").
				WillReturnRows(rows).
				WillReturnError(nil)

			ctx := context.WithValue(context.Background(), config.USERIDCONTEXTKEY, "test")
			cards, err := ds.GetCards(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetCards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(cards) != tt.cardCount {
				t.Errorf("len(GetCards()) = %v, want %v", len(cards), tt.cardCount)
			}
		})
	}
}

func TestDBStorage_CreateCard(t *testing.T) {
	timeNow := time.Now()
	tests := []struct {
		name    string
		card    models.Card
		setup   func(mock sqlmock.Sqlmock)
		want    models.Card
		wantErr bool
	}{
		{
			name: "Valid card",
			card: models.NewCard("1234 5678 9012 3456", "12/24", "Test User", "123"),
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO cards(number, expiration_date, holder_name, cvv, user_id) VALUES ($1, $2, $3, $4, $5)`)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, uploaded_at FROM cards WHERE number=$1 and user_id=$2`)).WillReturnRows(sqlmock.NewRows([]string{"id", "uploaded_at"}).AddRow("1", timeNow))
			},
			want:    models.Card{ID: "1", Number: "1234 5678 9012 3456", ExpirationDate: "12/24", HolderName: "Test User", CVV: "123", UploadedAt: timeNow},
			wantErr: false,
		},
		{
			name: "Card with same number already exists",
			card: models.NewCard("9876 5432 1098 7654", "07/26", "User Test", "321"),
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO cards(number, expiration_date, holder_name, cvv, user_id) VALUES ($1, $2, $3, $4, $5)`)).WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})
			},
			want:    models.NewCard("9876 5432 1098 7654", "07/26", "User Test", "321"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			ds := &DBStorage{
				db: db,
			}

			tt.setup(mock)

			ctx := context.WithValue(context.Background(), config.USERIDCONTEXTKEY, "test")
			card, err := ds.CreateCard(ctx, tt.card)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(card, tt.want) {
				t.Errorf("CreateCard() = %v, want %v", card, tt.want)
			}
		})
	}
}
func TestDBStorage_DeleteCredentials(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*DBStorage, sqlmock.Sqlmock, string)
		wantErr bool
	}{
		{
			name: "Valid Credentials ID",
			setup: func(ds *DBStorage, mock sqlmock.Sqlmock, userID string) {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM credentials WHERE user_id=$1 and id=$2`)).
					WithArgs(userID, "1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Query Execution Error",
			setup: func(ds *DBStorage, mock sqlmock.Sqlmock, userID string) {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM credentials WHERE user_id=$1 and id=$2`)).
					WithArgs(userID, "9999").
					WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}
	ctx := context.WithValue(context.Background(), config.USERIDCONTEXTKEY, "test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			ds := &DBStorage{db: db}
			tt.setup(ds, mock, ctx.Value(config.USERIDCONTEXTKEY).(string))
			err := ds.DeleteCredentials(ctx, "1")
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDBStorage_UpdateCredentials(t *testing.T) {
	timeNow := time.Now()
	tests := []struct {
		name        string
		credentials models.Credentials
		setup       func(*DBStorage, sqlmock.Sqlmock, string)
		want        models.Credentials
		wantErr     bool
	}{
		{
			name: "Valid Update",
			credentials: models.Credentials{
				ID:          "1",
				ServiceName: "Gmail",
				Identity:    "user@gmail.com",
				Password:    "password123",
				UploadedAt:  timeNow,
			},
			setup: func(ds *DBStorage, mock sqlmock.Sqlmock, userID string) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE credentials SET service_name=$1, identity=$2, password=$3 WHERE user_id=$4 and id=$5`)).
					WithArgs("Gmail", "user@gmail.com", "password123", userID, "1").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT uploaded_at FROM credentials WHERE id=$1 and user_id=$2`)).
					WithArgs("1", userID).
					WillReturnRows(sqlmock.NewRows([]string{"uploaded_at"}).AddRow(timeNow))
			},
			want: models.Credentials{
				ID:          "1",
				ServiceName: "Gmail",
				Identity:    "user@gmail.com",
				Password:    "password123",
				UploadedAt:  timeNow,
			},
			wantErr: false,
		},
		{
			name: "Invalid Update",
			credentials: models.Credentials{
				ID:          "9999",
				ServiceName: "Yahoo",
				Identity:    "user@yahoo.com",
				Password:    "password456",
				UploadedAt:  time.Now(),
			},
			setup: func(ds *DBStorage, mock sqlmock.Sqlmock, userID string) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE credentials SET service_name=$1, identity=$2, password=$3 WHERE user_id=$4 and id=$5`)).
					WithArgs("Yahoo", "user@yahoo.com", "password456", userID, "9999").
					WillReturnError(fmt.Errorf("some error"))
			},
			want:    models.Credentials{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			ds := &DBStorage{db: db}
			ctx := context.WithValue(context.Background(), config.USERIDCONTEXTKEY, "test")
			tt.setup(ds, mock, ctx.Value(config.USERIDCONTEXTKEY).(string))
			got, err := ds.UpdateCredentials(ctx, tt.credentials)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateCredentials() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBStorage_GetCredentials(t *testing.T) {
	timeNow := time.Now()
	tests := []struct {
		name    string
		setup   func(*DBStorage, sqlmock.Sqlmock, string)
		want    []models.Credentials
		wantErr bool
	}{
		{
			name: "Valid Get",
			setup: func(ds *DBStorage, mock sqlmock.Sqlmock, userID string) {
				mockRows := sqlmock.NewRows([]string{"id", "service_name", "identity", "password", "uploaded_at"}).
					AddRow("1", "Service", "Identity", "Password", timeNow)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, service_name, identity, password, uploaded_at FROM credentials WHERE user_id=$1`)).
					WithArgs("test").
					WillReturnRows(mockRows)
			},
			want: []models.Credentials{
				{
					ID:          "1",
					ServiceName: "Service",
					Identity:    "Identity",
					Password:    "Password",
					UploadedAt:  timeNow,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Get",
			setup: func(ds *DBStorage, mock sqlmock.Sqlmock, userID string) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, service_name, identity, password, uploaded_at FROM credentials WHERE user_id=$1`)).
					WithArgs("invalid").
					WillReturnError(fmt.Errorf("some error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			ds := &DBStorage{db: db}
			ctx := context.WithValue(context.Background(), config.USERIDCONTEXTKEY, "test")
			tt.setup(ds, mock, ctx.Value(config.USERIDCONTEXTKEY).(string))
			got, err := ds.GetCredentials(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCredentials() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBStorage_CreateCredentials(t *testing.T) {
	timeNow := time.Now()
	tests := []struct {
		name        string
		credentials models.Credentials
		setup       func(*DBStorage, sqlmock.Sqlmock, string)
		want        models.Credentials
		wantErr     bool
	}{
		{
			name: "Valid Insert",
			credentials: models.Credentials{
				ServiceName: "Facebook",
				Identity:    "tester@facebook.com",
				Password:    "testpassword",
			},
			setup: func(ds *DBStorage, mock sqlmock.Sqlmock, userID string) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO credentials(service_name, identity, password, user_id) VALUES ($1, $2, $3, $4)`)).
					WithArgs("Facebook", "tester@facebook.com", "testpassword", userID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, uploaded_at FROM credentials WHERE service_name=$1 and user_id=$2`)).
					WithArgs("Facebook", userID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "uploaded_at"}).AddRow("1", timeNow))
			},
			want: models.Credentials{
				ID:          "1",
				ServiceName: "Facebook",
				Identity:    "tester@facebook.com",
				Password:    "testpassword",
				UploadedAt:  timeNow,
			},
			wantErr: false,
		},
		{
			name: "Invalid Insert",
			credentials: models.Credentials{
				ServiceName: "Twitter",
				Identity:    "tester@twitter.com",
				Password:    "testpassword",
			},
			setup: func(ds *DBStorage, mock sqlmock.Sqlmock, userID string) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO credentials(service_name, identity, password, user_id) VALUES ($1, $2, $3, $4)`)).
					WithArgs("Twitter", "tester@twitter.com", "testpassword", userID).
					WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})
			},
			want: models.Credentials{
				ServiceName: "Twitter",
				Identity:    "tester@twitter.com",
				Password:    "testpassword",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			ds := &DBStorage{db: db}
			ctx := context.WithValue(context.Background(), config.USERIDCONTEXTKEY, "test")
			tt.setup(ds, mock, ctx.Value(config.USERIDCONTEXTKEY).(string))
			got, err := ds.CreateCredentials(ctx, tt.credentials)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateCredentials() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBStorage_AuthorizeUser(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		setup   func(*DBStorage, sqlmock.Sqlmock, string)
		want    models.User
		wantErr bool
	}{
		{
			name:  "Valid User",
			email: "test@test.com",
			setup: func(ds *DBStorage, mock sqlmock.Sqlmock, email string) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password FROM users WHERE email=$1`)).
					WithArgs(email).WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).AddRow("1", "TestPassword"))
			},
			want: models.User{
				ID:       "1",
				Email:    "test@test.com",
				Password: "TestPassword",
			},
			wantErr: false,
		},
		{
			name:  "Invalid User",
			email: "invalid@test.com",
			setup: func(ds *DBStorage, mock sqlmock.Sqlmock, email string) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password FROM users WHERE email=$1`)).
					WithArgs(email).WillReturnError(sql.ErrNoRows)
			},
			want:    models.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			ds := &DBStorage{db: db}
			ctx := context.Background()
			tt.setup(ds, mock, tt.email)
			got, err := ds.AuthorizeUser(ctx, tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthorizeUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthorizeUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBStorage_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    models.User
		setup   func(*DBStorage, sqlmock.Sqlmock, string)
		want    models.User
		wantErr bool
	}{
		{
			name: "Valid New User",
			user: models.User{
				Email:    "newvaliduser@test.com",
				Password: "testpassword",
			},
			setup: func(ds *DBStorage, mock sqlmock.Sqlmock, email string) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users(email, password) VALUES ($1, $2)`)).
					WithArgs(email, "testpassword").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password FROM users WHERE email=$1`)).
					WithArgs(email).
					WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).
						AddRow("1", "testpassword"))
			},
			want: models.User{
				ID:       "1",
				Email:    "newvaliduser@test.com",
				Password: "testpassword",
			},
			wantErr: false,
		},
		{
			name: "Valid Existing User",
			user: models.User{
				Email:    "existinguser@test.com",
				Password: "testpassword",
			},
			setup: func(ds *DBStorage, mock sqlmock.Sqlmock, email string) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users(email, password) VALUES ($1, $2)`)).
					WithArgs(email, "testpassword").
					WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})
			},
			want:    models.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			ds := &DBStorage{db: db}
			ctx := context.Background()
			tt.setup(ds, mock, tt.user.Email)
			got, err := ds.CreateUser(ctx, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
