package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

	// An interface to interact with repository
	DatabaseInteractor interface {
		Database() Storage
		CloseConnection() error
		DDCredentials() DBCredentials
	}

	databaseInteractor struct {
		DB     Storage
		DBCred DBCredentials
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

	var interactor databaseInteractor
	driver := cred.GetDriver()
	if driver == "postgres" {
		interactor.DB = newPGStorage(database)
	} else if driver == "mysql" {
		interactor.DB = newMySQLStorage(database)
	} else {
		return nil, fmt.Errorf("error invalid driver %s", driver)
	}

	interactor.DBCred = cred
	return &interactor, nil
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
			return nil, fmt.Errorf("error connecting to database: %w", handleError(err))
		}
	case "mysql":
		connStr = fmt.Sprintf("mysql://%s:%s@%s/%s",
			c.GetUser(), c.GetSecret(), c.GetAddress(), c.GetDBName())

		db, err = sql.Open(c.GetDriver(), connStr)
		if err != nil {
			return nil, fmt.Errorf("error connecting to database: %w", handleError(err))
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

func (d *databaseInteractor) Database() Storage {
	return d.DB
}

func (d *databaseInteractor) DDCredentials() DBCredentials {
	return d.DBCred
}

func (d *databaseInteractor) CloseConnection() error {
	var err error

	ctx, timeoutCloser := context.WithTimeout(context.Background(), 15*time.Second)
	defer timeoutCloser()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		time.Sleep(30 * time.Second)
		err = fmt.Errorf("error closing connection to database. Aborting closer due to timeout")
		cancel()
	}()

	err = d.DB.Close()

	return err
}
