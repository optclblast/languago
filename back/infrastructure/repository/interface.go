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
		Txer
		CreateUser(ctx context.Context, arg CreateUserParams) error
		UpdateUser(ctx context.Context, arg UpdateUserParams) error
		DeleteUser(ctx context.Context, userID uuid.UUID) error
		SelectUser(ctx context.Context, arg SelectUserParams) (*entities.User, error)
	}

	FlashcardRepository interface {
		Txer
		CreateFlashcard(ctx context.Context, arg CreateFlashcardParams) error
		UpdateFlashcard(ctx context.Context, arg UpdateFlashcardParams) error
		DeleteFlashcard(ctx context.Context, cardID uuid.UUID) error
		SelectFlashcard(ctx context.Context, arg SelectFlashcardParams) ([]*entities.Flashcard, error)
	}

	DeckRepository interface {
		Txer
		CreateDeck(ctx context.Context, arg CreateDeckParams) error
		UpdateDeck(ctx context.Context, arg UpdateDeckParams) error
		DeleteDeck(ctx context.Context, deckID uuid.UUID) error
		SelectDeck(ctx context.Context, arg SelectDeckParams) (*entities.Deck, error)
		AddToDeck(ctx context.Context, arg AddToDeckParams) error
		DeleteFromDeck(ctx context.Context, arg DeleteFromDeckParams) error
		SelectFromDeck(ctx context.Context, arg SelectFromDeckParams) (*entities.Flashcard, error)
	}

	// Storage interface provides an abstraction over particular database used by node
	Storage interface {
		PingDB() error
		Close() error

		UserRepository
		FlashcardRepository
		DeckRepository
	}

	Txer interface {
		WithTransaction(context.Context, *sql.Tx, func(*sql.Tx) error) error
	}

	pgStorage struct {
		conn *sql.DB
		db   *postgresql.Queries
	}
)

// Storage implementation for PostgreSQL database
func newPGStorage(db *sql.DB) *pgStorage {
	return &pgStorage{
		conn: db,
		db:   postgresql.New(db),
	}
}

func (s *pgStorage) PingDB() error {
	if err := s.conn.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}
	return nil
}

func (s *pgStorage) Close() error {
	return s.conn.Close()
}

func (s *pgStorage) WithTransaction(ctx context.Context, tx *sql.Tx, txFunc func(*sql.Tx) error) error {
	var err error

	defer func() {
		// context deadline exceeded - try to rollback transaction
		if _, ok := <-ctx.Done(); !ok && tx != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				err = errors.Join(err, rbErr)
			}
		}
	}()

	var hasExternalTx bool = true
	if tx == nil {
		isolationLevel := func() sql.IsolationLevel {
			if isolationLevel, err := ctxtools.IsolationLevel(ctx); err == nil {
				return isolationLevel
			}

			return globalIsolationLevel
		}()

		tx, err = s.conn.BeginTx(ctx, &sql.TxOptions{
			Isolation: isolationLevel,
		})
		if err != nil {
			return fmt.Errorf("error begin transaction: %w", err)
		}

		hasExternalTx = false
	}

	err = txFunc(tx)
	if err != nil {
		return fmt.Errorf("error run transaction function: %w", err)
	}

	if !hasExternalTx {
		err = tx.Commit()
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("error rollback transaction: %w", rbErr)
			}
			return fmt.Errorf("error commit transaction: %w", err)
		}
	}

	return nil
}

func (s *pgStorage) CreateUser(ctx context.Context, arg CreateUserParams) error {
	// todo move validation from repository
	if len(arg.Login) < 4 && arg.Password == "" {
		return fmt.Errorf("error invalid user credentials: %w", ErrInvalidData)
	}

	err := s.db.CreateUser(ctx, postgresql.CreateUserParams{
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
		_, err := s.db.UpdateUserLogin(ctx, postgresql.UpdateUserLoginParams{
			ID:    arg.ID,
			Login: sql.NullString{String: arg.Login, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("error update user login: %w", err)
		}
	}

	if arg.Password != "" {
		err := s.db.UpdateUserPassword(ctx, postgresql.UpdateUserPasswordParams{
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

	err := s.db.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("error delete user: %w", err)
	}

	return nil
}

func (s *pgStorage) SelectUser(ctx context.Context, arg SelectUserParams) (*entities.User, error) {
	var user postgresql.User
	var err error

	err = s.WithTransaction(ctx, nil, func(tx *sql.Tx) error {
		err = s.db.AddToDeck(ctx, postgresql.AddToDeckParams{})
		err = s.db.CreateUser(ctx, postgresql.CreateUserParams{})

		return err
	})

	if arg.ID != uuid.Nil {
		user, err = s.db.SelectUserByID(ctx, arg.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errors2.ErrNotFound
			}
		}
	} else if arg.Login != "" {
		user, err = s.db.SelectUserByLogin(ctx, sql.NullString{String: arg.Login, Valid: true})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errors2.ErrNotFound
			}
		}
	}

	return entities.UserFromPG(user), nil
}

func (s *pgStorage) CreateFlashcard(ctx context.Context, args CreateFlashcardParams) error {
	_, err := s.db.CreateFlashcard(ctx, postgresql.CreateFlashcardParams{
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

	currentFlashcardState, err := s.db.SelectFlashcardByID(ctx, arg.ID)
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

	err = s.db.UpdateFlashcard(ctx, *newVals)
	if err != nil {
		return fmt.Errorf("error updating flashcard: %w", handleError(err))
	}

	return nil
}
func (s *pgStorage) DeleteFlashcard(ctx context.Context, cardID uuid.UUID) error {
	if cardID == uuid.Nil {
		return fmt.Errorf("error flashcard uuid is required")
	}

	err := s.db.DeleteFlashcard(ctx, cardID)
	if err != nil {
		return fmt.Errorf("error delete flashcard: %w", err)
	}

	return nil
}

func (s *pgStorage) SelectFlashcard(ctx context.Context, arg SelectFlashcardParams) ([]*entities.Flashcard, error) {
	// switch {
	// case arg.ID != uuid.Nil:
	// 	card, err := s.db.SelectFlashcardByID(ctx, arg.ID)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	flashcard, err :=
	// }
	return nil, nil
}

func (s *pgStorage) CreateDeck(ctx context.Context, arg CreateDeckParams) error {
	_, err := s.db.CreateDeck(ctx, postgresql.CreateDeckParams{
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
	err := s.db.AddToDeck(ctx, postgresql.AddToDeckParams{
		DeckID:      uuid.NullUUID{UUID: arg.DeckID, Valid: true},
		FlashcardID: uuid.NullUUID{UUID: arg.FlashcardID, Valid: true},
	})
	return err
}

func (s *pgStorage) DeleteFromDeck(ctx context.Context, arg DeleteFromDeckParams) error { return nil }

func (s *pgStorage) SelectFromDeck(ctx context.Context, arg SelectFromDeckParams) (*entities.Flashcard, error) {
	return nil, nil
}
