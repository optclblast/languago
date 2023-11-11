package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	errors2 "languago/internal/pkg/errors"
	"languago/internal/pkg/models/entities"
	"languago/internal/pkg/repository/postgresql"

	"github.com/google/uuid"
)

type (
	UserRepository interface {
		CreateUser(ctx context.Context, arg CreateUserParams) error
		UpdateUser(ctx context.Context, arg UpdateUserParams) error
		DeleteUser(ctx context.Context, userID uuid.UUID) error
		SelectUser(ctx context.Context, arg SelectUserParams) (*entities.User, error)
	}

	// Storage interface provides an abstraction over particular database used by node
	Storage interface {
		// Ping func
		PingDB() error
		Close() error

		CreateUser(ctx context.Context, arg CreateUserParams) error
		UpdateUser(ctx context.Context, arg UpdateUserParams) error
		DeleteUser(ctx context.Context, userID uuid.UUID) error
		SelectUser(ctx context.Context, arg SelectUserParams) (*entities.User, error)

		CreateFlashcard(ctx context.Context, arg CreateFlashcardParams) error
		UpdateFlashcard(ctx context.Context, arg UpdateFlashcardParams) error
		DeleteFlashcard(ctx context.Context, cardID uuid.UUID) error
		SelectFlashcard(ctx context.Context, arg SelectFlashcardParams) ([]*entities.Flashcard, error)

		CreateDeck(ctx context.Context, arg CreateDeckParams) error
		UpdateDeck(ctx context.Context, arg UpdateDeckParams) error
		DeleteDeck(ctx context.Context, deckID uuid.UUID) error
		SelectDeck(ctx context.Context, arg SelectDeckParams) (*entities.Deck, error)

		AddToDeck(ctx context.Context, arg AddToDeckParams) error
		DeleteFromDeck(ctx context.Context, arg DeleteFromDeckParams) error
		SelectFromDeck(ctx context.Context, arg SelectFromDeckParams) (*entities.Flashcard, error)
	}

	pgStorage struct {
		conn *sql.DB
		db   *postgresql.Queries
	}

	mysqlStorage struct {
		conn *sql.DB
		//db *mysql.Queries
	}
)

// Storage implementation for PostgreSQL database
func newPGStorage(db *sql.DB) *pgStorage {
	return &pgStorage{
		conn: db,
		db:   postgresql.New(db),
	}
}

func newMySQLStorage(db *sql.DB) *mysqlStorage {
	return nil
	// return &mysqlStorage{
	// 	db: mysql.New(db),
	// }
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

func (s *pgStorage) CreateUser(ctx context.Context, arg CreateUserParams) error {
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

func (s *pgStorage) CreateFlashcard(ctx context.Context, arg CreateFlashcardParams) error { return nil }

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

	return nil, nil
}

func (s *pgStorage) CreateDeck(ctx context.Context, arg CreateDeckParams) error { return nil }
func (s *pgStorage) UpdateDeck(ctx context.Context, arg UpdateDeckParams) error {
	return nil
}
func (s *pgStorage) DeleteDeck(ctx context.Context, deckID uuid.UUID) error { return nil }
func (s *pgStorage) SelectDeck(ctx context.Context, arg SelectDeckParams) (*entities.Deck, error) {
	return nil, nil
}

func (s *pgStorage) AddToDeck(ctx context.Context, arg AddToDeckParams) error           { return nil }
func (s *pgStorage) DeleteFromDeck(ctx context.Context, arg DeleteFromDeckParams) error { return nil }
func (s *pgStorage) SelectFromDeck(ctx context.Context, arg SelectFromDeckParams) (*entities.Flashcard, error) {
	return nil, nil
}

// Storage implementation for MySQL database
func (s *mysqlStorage) PingDB() error {
	if err := s.conn.Ping(); err != nil {
		s.conn.Close()
		return fmt.Errorf("error pinging database: %w", err)
	}
	return nil
}

func (s *mysqlStorage) Close() error {
	return s.conn.Close()
}

func (s *mysqlStorage) CreateUser(ctx context.Context, arg CreateUserParams) error { return nil }
func (s *mysqlStorage) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	return nil
}
func (s *mysqlStorage) DeleteUser(ctx context.Context, userID uuid.UUID) error { return nil }
func (s *mysqlStorage) SelectUser(ctx context.Context, arg SelectUserParams) (*entities.User, error) {
	return nil, nil
}

func (s *mysqlStorage) CreateFlashcard(ctx context.Context, arg CreateFlashcardParams) error {
	return nil
}

// Updated the flashcard record. ID must be not nil.
func (s *mysqlStorage) UpdateFlashcard(ctx context.Context, arg UpdateFlashcardParams) error {
	// if arg.ID == uuid.Nil {
	// 	return fmt.Errorf("error id required")
	// }

	// currentFlashcardState, err := s.db.SelectFlashcardByID(ctx, arg.ID)
	// if err != nil {
	// 	return fmt.Errorf("error selecting flashcard: %w", err)
	// }

	// newVals := &postgresql.UpdateFlashcardParams{
	// 	Word:    currentFlashcardState.Word,
	// 	Meaning: currentFlashcardState.Meaning,
	// 	Usage:   currentFlashcardState.Usage,
	// 	ID:      currentFlashcardState.ID,
	// }

	// switch {
	// case arg.Meaning != "":
	// 	newVals.Meaning = sql.NullString{String: arg.Meaning, Valid: true}
	// case arg.Word != "":
	// 	newVals.Word = sql.NullString{String: arg.Word, Valid: true}
	// case arg.Usage != nil:
	// 	newVals.Usage = arg.Usage
	// default:
	// 	return nil
	// }

	// err = s.db.UpdateFlashcard(ctx, *newVals)
	// if err != nil {
	// 	return fmt.Errorf("error updating flashcard: %w", err)
	// }

	return nil
}
func (s *mysqlStorage) DeleteFlashcard(ctx context.Context, cardID uuid.UUID) error { return nil }
func (s *mysqlStorage) SelectFlashcard(ctx context.Context, arg SelectFlashcardParams) ([]*entities.Flashcard, error) {
	return nil, nil
}

func (s *mysqlStorage) CreateDeck(ctx context.Context, arg CreateDeckParams) error { return nil }
func (s *mysqlStorage) UpdateDeck(ctx context.Context, arg UpdateDeckParams) error {
	return nil
}
func (s *mysqlStorage) DeleteDeck(ctx context.Context, deckID uuid.UUID) error { return nil }
func (s *mysqlStorage) SelectDeck(ctx context.Context, arg SelectDeckParams) (*entities.Deck, error) {
	return nil, nil
}

func (s *mysqlStorage) AddToDeck(ctx context.Context, arg AddToDeckParams) error { return nil }
func (s *mysqlStorage) DeleteFromDeck(ctx context.Context, arg DeleteFromDeckParams) error {
	return nil
}
func (s *mysqlStorage) SelectFromDeck(ctx context.Context, arg SelectFromDeckParams) (*entities.Flashcard, error) {
	return nil, nil
}
