package repository

import (
	"context"
	"database/sql"
	"fmt"
	"languago/internal/pkg/models/requests/rest"
	"languago/internal/pkg/repository/postgresql"

	"github.com/google/uuid"
)

func (d *databaseInteractor) EditFlashcard(ctx context.Context, arg *rest.EditFlashcardRequest) error {
	uuid, err := uuid.Parse(arg.Id)
	if err != nil {
		return fmt.Errorf("error invalid id: %w", err)
	}

	card, err := d.Database().SelectFlashcardByID(ctx, uuid)
	if err != nil {
		return fmt.Errorf("error selecting flashcard: %w", err)
	}

	newVals := &postgresql.UpdateFlashcardParams{
		Word:    card.Word,
		Meaning: card.Meaning,
		Usage:   card.Usage,
		ID:      uuid,
	}

	switch {
	case arg.WordInNative != "":
		newVals.Meaning = sql.NullString{String: arg.WordInNative, Valid: true}
	case arg.WordInTarget != "":
		newVals.Word = sql.NullString{String: arg.WordInTarget, Valid: true}
	case arg.UsageExamples != nil:
		newVals.Usage = arg.UsageExamples
	default:
		return nil
	}

	err = d.Database().UpdateFlashcard(ctx, *newVals)
	if err != nil {
		return fmt.Errorf("error updating flashcard: %w", err)
	}

	return nil
}
