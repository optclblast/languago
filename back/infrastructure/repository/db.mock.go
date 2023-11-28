package repository

import (
	"context"
	"languago/pkg/models/entities"
	"languago/test/generators"
	"math/rand"

	"github.com/google/uuid"
)

type mockStorage struct{}

func _newMockStorage() Storage {
	return &mockStorage{}
}

func (s *mockStorage) PingDB() error {
	return nil
}

func (s *mockStorage) Close() error {
	return nil
}

func (s *mockStorage) CreateUser(ctx context.Context, arg CreateUserParams) error { return nil }
func (s *mockStorage) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	return nil
}
func (s *mockStorage) DeleteUser(ctx context.Context, userID uuid.UUID) error { return nil }
func (s *mockStorage) SelectUser(ctx context.Context, arg SelectUserParams) (*entities.User, error) {
	return nil, nil
}

func (s *mockStorage) CreateFlashcard(ctx context.Context, arg CreateFlashcardParams) error {
	return nil
}

func (s *mockStorage) UpdateFlashcard(ctx context.Context, arg UpdateFlashcardParams) error {
	return nil
}
func (s *mockStorage) DeleteFlashcard(ctx context.Context, cardID uuid.UUID) error { return nil }
func (s *mockStorage) SelectFlashcard(ctx context.Context, arg SelectFlashcardParams) ([]*entities.Flashcard, error) {
	len := rand.Intn(20)
	resp := make([]*entities.Flashcard, 0, len)

	for i := 0; i < len; i++ {
		card := entities.Flashcard{
			ID:             uuid.New(),
			NativeLanguage: generators.RandStringRunes(10),
			TargetLang:     generators.RandStringRunes(10),
			Meaning:        generators.RandStringRunes(10),
			Word:           generators.RandStringRunes(10),
			UsageExamples:  generators.RandStringSlice(5, 15),
		}
		resp = append(resp, &card)
	}

	return resp, nil
}

func (s *mockStorage) CreateDeck(ctx context.Context, arg CreateDeckParams) error { return nil }
func (s *mockStorage) UpdateDeck(ctx context.Context, arg UpdateDeckParams) error {
	return nil
}
func (s *mockStorage) DeleteDeck(ctx context.Context, deckID uuid.UUID) error { return nil }
func (s *mockStorage) SelectDeck(ctx context.Context, arg SelectDeckParams) (*entities.Deck, error) {
	return nil, nil
}

func (s *mockStorage) AddToDeck(ctx context.Context, arg AddToDeckParams) error { return nil }
func (s *mockStorage) DeleteFromDeck(ctx context.Context, arg DeleteFromDeckParams) error {
	return nil
}
func (s *mockStorage) SelectFromDeck(ctx context.Context, arg SelectFromDeckParams) (*entities.Flashcard, error) {
	return nil, nil
}
