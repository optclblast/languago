package repository

import (
	"github.com/google/uuid"
)

type (
	// PARAMS
	CreateUserParams struct {
		ID       uuid.UUID `db:"id" json:"id"`
		Login    string    `db:"login" json:"login"`
		Password string    `db:"password" json:"password"`
	}

	AddToDeckParams struct {
		DeckID      uuid.UUID `db:"deck_id" json:"deck_id"`
		FlashcardID uuid.UUID `db:"flashcard_id" json:"flashcard_id"`
	}

	CreateDeckParams struct {
		ID    uuid.UUID `db:"id" json:"id"`
		Name  string    `db:"name" json:"name"`
		Owner uuid.UUID `db:"owner" json:"owner"`
	}

	CreateFlashcardParams struct {
		ID      uuid.UUID `db:"id" json:"id"`
		Word    string    `db:"word" json:"word"`
		Meaning string    `db:"meaning" json:"meaning"`
		Usage   []string  `db:"usage" json:"usage"`
	}

	DeleteFromDeckParams struct {
		FlashcardID uuid.UUID `db:"flashcard_id" json:"flashcard_id"`
		DeckID      uuid.UUID `db:"deck_id" json:"deck_id"`
	}

	UpdateDeckParams struct {
		Name string    `db:"name" json:"name"`
		ID   uuid.UUID `db:"id" json:"id"`
	}

	SelectFlashcardParams struct {
		ID          uuid.UUID `db:"id" json:"id"`
		Word        string    `db:"word" json:"word"`
		Meaning     string    `db:"meaning" json:"meaning"`
		Usage       []string  `db:"usage" json:"usage"`
		DeckID      uuid.UUID `db:"deck_id" json:"deck_id"`
		FlashcardID uuid.UUID `db:"flashcard_id" json:"flashcard_id"`
	}

	SelectUserParams struct {
		ID    uuid.UUID `db:"id" json:"id"`
		Login string    `db:"login" json:"login"`
	}

	SelectDeckParams struct {
		ID    uuid.UUID `db:"id" json:"id"`
		Name  string    `db:"name" json:"name"`
		Owner uuid.UUID `db:"owner" json:"owner"`
	}

	UpdateFlashcardParams struct {
		Word    string    `db:"word" json:"word"`
		Meaning string    `db:"meaning" json:"meaning"`
		Usage   []string  `db:"usage" json:"usage"`
		ID      uuid.UUID `db:"id" json:"id"`
	}

	SelectFromDeckParams struct {
		CardID      uuid.UUID `db:"card_id" json:"card_id"`
		DeckID      uuid.UUID `db:"deck_id" json:"deck_id"`
		DeckOwner   uuid.UUID `db:"deck_owner" json:"deck_owner"`
		WordMeaning string    `db:"word_meaning" json:"word_meaning"`
		Word        string    `db:"word" json:"word"`
		Usage       []string  `db:"usage" json:"usage"`
	}

	UpdateUserParams struct {
		Login    string    `db:"login" json:"login"`
		ID       uuid.UUID `db:"id" json:"id"`
		Password string    `db:"password" json:"password"`
	}
)
