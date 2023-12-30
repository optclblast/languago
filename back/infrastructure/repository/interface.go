package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"languago/infrastructure/repository/postgresql"
	"languago/pkg/ctxtools"
	errors2 "languago/pkg/errors"
	"languago/pkg/models/entities"

	"github.com/google/uuid"
)

var globalIsolationLevel sql.IsolationLevel

type (
	UserRepository interface {
		TransactionalStorage
		CreateUser(ctx context.Context, arg CreateUserParams) error
		UpdateUser(ctx context.Context, arg UpdateUserParams) error
		DeleteUser(ctx context.Context, userID uuid.UUID) error
		SelectUser(ctx context.Context, arg SelectUserParams) (*entities.User, error)
	}

	FlashcardRepository interface {
		TransactionalStorage
		CreateFlashcard(ctx context.Context, arg CreateFlashcardParams) error
		UpdateFlashcard(ctx context.Context, arg UpdateFlashcardParams) error
		DeleteFlashcard(ctx context.Context, cardID uuid.UUID) error
		SelectFlashcard(ctx context.Context, arg SelectFlashcardParams) ([]*entities.Flashcard, error)
	}

	DeckRepository interface {
		TransactionalStorage
		CreateDeck(ctx context.Context, arg CreateDeckParams) error
		UpdateDeck(ctx context.Context, arg UpdateDeckParams) error
		DeleteDeck(ctx context.Context, deckID uuid.UUID) error
		SelectDeck(ctx context.Context, arg SelectDeckParams) (*entities.Deck, error)
		AddToDeck(ctx context.Context, arg AddToDeckParams) error
		DeleteFromDeck(ctx context.Context, arg DeleteFromDeckParams) error
		SelectFromDeck(ctx context.Context, arg SelectFromDeckParams) (*entities.Flashcard, error)
	}

	DBTX interface {
		ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
		PrepareContext(context.Context, string) (*sql.Stmt, error)
		QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
		QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	}

	Storage interface {
		PingDB() error
		Close() error

		UserRepository
		FlashcardRepository
		DeckRepository
	}

	TransactionalStorage interface {
		Conn(ctx context.Context) DBTX
		WithTransaction(context.Context, func(ctx context.Context) error) error
	}

	pgStorage struct {
		db      *sql.DB
		querier *postgresql.Queries
	}
)

// Storage implementation for PostgreSQL database
func newPGStorage(db *sql.DB) *pgStorage {
	return &pgStorage{
		db:      db,
		querier: postgresql.New(db),
	}
}

func (s *pgStorage) PingDB() error {
	if err := s.db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}
	return nil
}

func (s *pgStorage) Close() error {
	return s.db.Close()
}

type txCtxKey struct{}

func (s *pgStorage) Conn(ctx context.Context) DBTX {
	if tx, ok := ctx.Value(txCtxKey{}).(*sql.Tx); ok {
		return tx
	}

	return s.db
}

func (s *pgStorage) WithTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error {
	var err error
	var tx *sql.Tx = new(sql.Tx)

	defer func() {
		if err == nil {
			err = tx.Commit()
		}

		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("error rollback transaction: %w", rbErr)
				return
			}
			err = fmt.Errorf("error commit transaction: %w", err)
			return
		}
	}()

	if _, ok := ctx.Value(txCtxKey{}).(*sql.Tx); !ok {
		isolationLevel := func() sql.IsolationLevel {
			if isolationLevel, err := ctxtools.IsolationLevel(ctx); err == nil {
				return isolationLevel
			}

			return globalIsolationLevel
		}()

		tx, err = s.db.BeginTx(ctx, &sql.TxOptions{
			Isolation: isolationLevel,
		})
		if err != nil {
			return fmt.Errorf("error begin transaction: %w", err)
		}

		ctx = context.WithValue(ctx, txCtxKey{}, tx)
	}

	err = txFunc(ctx)
	if err != nil {
		return fmt.Errorf("error run transaction function: %w", err)
	}

	return err
}

func (s *pgStorage) CreateUser(ctx context.Context, arg CreateUserParams) error {
	// todo move validation from repository
	if len(arg.Login) < 4 && arg.Password == "" {
		return fmt.Errorf("error invalid user credentials: %w", ErrInvalidData)
	}

	err := s.querier.CreateUser(ctx, postgresql.CreateUserParams{
		ID:       arg.ID,
		Login:    sql.NullString{String: arg.Login, Valid: true},
		Password: sql.NullString{String: arg.Password, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("error create user: %w", err)
	}

	return nil
}

func (s *pgStorage) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	// todo move validation from repository
	if arg.ID == uuid.Nil {
		return fmt.Errorf("error user id is required")
	}

	if arg.Login != "" {
		_, err := s.querier.UpdateUserLogin(ctx, postgresql.UpdateUserLoginParams{
			ID:    arg.ID,
			Login: sql.NullString{String: arg.Login, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("error update user login: %w", err)
		}
	}

	if arg.Password != "" {
		err := s.querier.UpdateUserPassword(ctx, postgresql.UpdateUserPasswordParams{
			ID:       arg.ID,
			Password: sql.NullString{String: arg.Login, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("error update user login: %w", err)
		}
	}

	return nil
}

func (s *pgStorage) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// todo move validation from repository
	if userID == uuid.Nil {
		return fmt.Errorf("error user id is required")
	}

	err := s.querier.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("error delete user: %w", err)
	}

	return nil
}

func (s *pgStorage) SelectUser(ctx context.Context, arg SelectUserParams) (*entities.User, error) {
	var user postgresql.User
	var err error

	err = s.WithTransaction(ctx, func(ctx context.Context) error {
		err = s.querier.AddToDeck(ctx, postgresql.AddToDeckParams{})
		err = s.querier.CreateUser(ctx, postgresql.CreateUserParams{})

		return err
	})

	if arg.ID != uuid.Nil {
		user, err = s.querier.SelectUserByID(ctx, arg.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errors2.ErrNotFound
			}
		}
	} else if arg.Login != "" {
		user, err = s.querier.SelectUserByLogin(ctx, sql.NullString{String: arg.Login, Valid: true})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errors2.ErrNotFound
			}
		}
	}

	return entities.UserFromPG(user), nil
}

func (s *pgStorage) CreateFlashcard(ctx context.Context, args CreateFlashcardParams) error {
	_, err := s.querier.CreateFlashcard(ctx, postgresql.CreateFlashcardParams{
		ID:      args.ID,
		Word:    sql.NullString{String: args.Word, Valid: true},
		Meaning: sql.NullString{String: args.Meaning, Valid: true},
		Usage:   args.Usage,
	})

	return err
}

func (s *pgStorage) UpdateFlashcard(ctx context.Context, arg UpdateFlashcardParams) error {
	if arg.ID == uuid.Nil {
		return fmt.Errorf("error id required")
	}

	currentFlashcardState, err := s.querier.SelectFlashcardByID(ctx, arg.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors2.ErrNotFound
		}

		return fmt.Errorf("error selecting flashcard: %w", handleError(err))
	}

	newVals := &postgresql.UpdateFlashcardParams{
		Word:    currentFlashcardState.Word,
		Meaning: currentFlashcardState.Meaning,
		Usage:   currentFlashcardState.Usage,
		ID:      currentFlashcardState.ID,
	}

	switch {
	case arg.Meaning != "":
		newVals.Meaning = sql.NullString{String: arg.Meaning, Valid: true}
	case arg.Word != "":
		newVals.Word = sql.NullString{String: arg.Word, Valid: true}
	case arg.Usage != nil:
		newVals.Usage = arg.Usage
	default:
		return nil
	}

	err = s.querier.UpdateFlashcard(ctx, *newVals)
	if err != nil {
		return fmt.Errorf("error updating flashcard: %w", handleError(err))
	}

	return nil
}
func (s *pgStorage) DeleteFlashcard(ctx context.Context, cardID uuid.UUID) error {
	if cardID == uuid.Nil {
		return fmt.Errorf("error flashcard uuid is required")
	}

	err := s.querier.DeleteFlashcard(ctx, cardID)
	if err != nil {
		return fmt.Errorf("error delete flashcard: %w", err)
	}

	return nil
}

func (s *pgStorage) SelectFlashcard(ctx context.Context, arg SelectFlashcardParams) ([]*entities.Flashcard, error) {
	// switch {
	// case arg.ID != uuid.Nil:
	// 	card, err := s.querier.SelectFlashcardByID(ctx, arg.ID)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	flashcard, err :=
	// }
	return nil, nil
}

func (s *pgStorage) CreateDeck(ctx context.Context, arg CreateDeckParams) error {
	_, err := s.querier.CreateDeck(ctx, postgresql.CreateDeckParams{
		ID:    arg.ID,
		Owner: uuid.NullUUID{UUID: arg.Owner, Valid: true},
		Name:  sql.NullString{String: arg.Name, Valid: true},
	})

	return err
}

func (s *pgStorage) UpdateDeck(ctx context.Context, arg UpdateDeckParams) error {
	return nil
}

func (s *pgStorage) DeleteDeck(ctx context.Context, deckID uuid.UUID) error { return nil }

func (s *pgStorage) SelectDeck(ctx context.Context, arg SelectDeckParams) (*entities.Deck, error) {
	return nil, nil
}

func (s *pgStorage) AddToDeck(ctx context.Context, arg AddToDeckParams) error {
	err := s.querier.AddToDeck(ctx, postgresql.AddToDeckParams{
		DeckID:      uuid.NullUUID{UUID: arg.DeckID, Valid: true},
		FlashcardID: uuid.NullUUID{UUID: arg.FlashcardID, Valid: true},
	})
	return err
}

func (s *pgStorage) DeleteFromDeck(ctx context.Context, arg DeleteFromDeckParams) error { return nil }

func (s *pgStorage) SelectFromDeck(ctx context.Context, arg SelectFromDeckParams) (*entities.Flashcard, error) {
	return nil, nil
}
