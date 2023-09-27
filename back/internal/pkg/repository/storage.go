package repository

import (
	"context"
	"database/sql"
	"fmt"
	"languago/internal/pkg/models/requests/rest"
	"languago/internal/pkg/repository/postgresql"

	_ "github.com/lib/pq"
)

type (
	DBCredentials interface {
		GetAddress() string
		GetDriver() string
		GetUser() string
		GetSecret() string
		GetDBName() string
		SetSSLMode(b bool)
		GetSSLMode() string
	}

	DBCred struct {
		DbAddress string
		DBName    string
		SSLMode   string
		Driver    string
		User      string
		Secret    string
	}

	DatabaseInteractor interface {
		Database() *postgresql.Queries
		DDCredentials() DBCredentials
		EditFlashcard(ctx context.Context, arg *rest.EditFlashcardRequest) error
	}

	databaseInteractor struct {
		DB     *postgresql.Queries
		DBCred DBCredentials
	}

	DBTX interface {
		ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
		PrepareContext(context.Context, string) (*sql.Stmt, error)
		QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
		QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	}
)

func NewDatabaseInteractor(cfg abstractDatabaseConfig) (DatabaseInteractor, error) {
	if cfg == nil {
		return nil, fmt.Errorf("error database config required.")
	}
	cred := cfg.GetCredentials()

	database, err := databaseConnection(cred)
	if err != nil {
		return nil, fmt.Errorf("error initializing database interactor: %w", err)
	}

	q := postgresql.New(database)

	return &databaseInteractor{
		DB:     q,
		DBCred: cred,
	}, nil
}

// TODO
func databaseConnection(c DBCredentials) (*sql.DB, error) {
	var (
		connStr string
		db      *sql.DB
		err     error
	)
	switch c.GetDriver() {
	case "postgres":
		connStr = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s",
			c.GetUser(), c.GetSecret(), c.GetAddress(), c.GetDBName(), c.GetSSLMode())

		db, err = sql.Open(c.GetDriver(), connStr)
		if err != nil {
			return nil, fmt.Errorf("error connecting to database: %w", err)
		}
	case "mysql":
		connStr = fmt.Sprintf("mysql://%s:%s@%s/%s",
			c.GetUser(), c.GetSecret(), c.GetAddress(), c.GetDBName())

		db, err = sql.Open(c.GetDriver(), connStr)
		if err != nil {
			return nil, fmt.Errorf("error connecting to database: %w", err)
		}
	default:
		return nil, fmt.Errorf("error connecting to database. unknown driver.")
	}

	return db, nil
}

// Relational database impl
func (c *DBCred) GetAddress() string {
	return c.DbAddress
}

func (c *DBCred) GetDriver() string {
	return c.Driver
}

func (c *DBCred) GetUser() string {
	return c.User
}

func (c *DBCred) GetSecret() string {
	return c.Secret
}

func (c *DBCred) GetDBName() string {
	return c.DBName
}

func (c *DBCred) SetSSLMode(b bool) {
	if b {
		c.SSLMode = "enabled"
	} else {
		c.SSLMode = "disabled"
	}
}

func (c *DBCred) GetSSLMode() string {
	return c.SSLMode
}

func (d *databaseInteractor) Database() *postgresql.Queries {
	return d.DB
}

func (d *databaseInteractor) DDCredentials() DBCredentials {
	return d.DBCred
}
