package pgsql

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Querier interface {
	AddToDeck(ctx context.Context, arg AddToDeckParams) error
	// Decks
	CreateDeck(ctx context.Context, arg CreateDeckParams) (CreateDeckRow, error)
	// Flashcards
	CreateFlashcard(ctx context.Context, arg CreateFlashcardParams) (CreateFlashcardRow, error)
	// User
	CreateUser(ctx context.Context, arg CreateUserParams) error
	DeleteDeck(ctx context.Context, id uuid.UUID) error
	DeleteFlashcard(ctx context.Context, id uuid.UUID) error
	DeleteFromDeck(ctx context.Context, arg DeleteFromDeckParams) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	EditDeckProps(ctx context.Context, arg EditDeckPropsParams) error
	SelectDeck(ctx context.Context, id uuid.UUID) (Deck, error)
	SelectDecksByName(ctx context.Context, name sql.NullString) ([]Deck, error)
	SelectFlashcardByID(ctx context.Context, id uuid.UUID) (Flashcard, error)
	SelectFlashcardByMeaning(ctx context.Context, arg SelectFlashcardByMeaningParams) ([]SelectFlashcardByMeaningRow, error)
	SelectFlashcardByWord(ctx context.Context, arg SelectFlashcardByWordParams) ([]SelectFlashcardByWordRow, error)
	SelectOwnerDecks(ctx context.Context, owner uuid.NullUUID) ([]Deck, error)
	SelectUser(ctx context.Context, arg SelectUserParams) (User, error)
	SelectUserByID(ctx context.Context, id uuid.UUID) (User, error)
	SelectUserByLogin(ctx context.Context, login sql.NullString) (User, error)
	UpdateFlashcard(ctx context.Context, arg UpdateFlashcardParams) error
	UpdateUserLogin(ctx context.Context, arg UpdateUserLoginParams) (UpdateUserLoginRow, error)
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error
}

var _ Querier = (*Queries)(nil)
