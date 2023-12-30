package flashcards

import (
	"context"
	"fmt"
	"languago/infrastructure/repository"
	"languago/pkg/models/requests/rest"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type FlashcardsController interface {
	CreateFlashcard(ctx context.Context, req *rest.NewFlashcardRequest) error
	GetFlashcard(ctx context.Context, args GetFlashcardParams) (*rest.GetFlashcardResponse, error)
	DeleteFlashcard(ctx context.Context, args DeleteFlashcardRequest) error
	EditFlashcard(ctx context.Context, args *rest.EditFlashcardRequest) error
}

type flashcardController struct {
	log     *logrus.Logger
	storage repository.DatabaseInteractor
}

func NewFlashcardsController(
	log *logrus.Logger,
	storage repository.DatabaseInteractor,
) FlashcardsController {
	return &flashcardController{
		log:     log,
		storage: storage,
	}
}

func (c *flashcardController) CreateFlashcard(ctx context.Context, req *rest.NewFlashcardRequest) error {
	err := c.storage.Database().CreateFlashcard(ctx, repository.CreateFlashcardParams{
		ID:      uuid.New(),
		Word:    req.Content.WordInTarget,
		Meaning: req.Content.WordInNative,
		Usage:   req.Content.UsageExamples,
	})
	if err != nil {
		return fmt.Errorf("error create flashcard: %w", err)
	}

	return nil
}

type GetFlashcardParams struct {
	Id      uuid.UUID
	DeckId  uuid.UUID
	Word    string
	Meaning string
}

func (c *flashcardController) GetFlashcard(ctx context.Context, args GetFlashcardParams) (*rest.GetFlashcardResponse, error) {
	return nil, nil
}

type DeleteFlashcardRequest struct{}

func (c *flashcardController) DeleteFlashcard(ctx context.Context, args DeleteFlashcardRequest) error {
	return nil
}

func (c *flashcardController) EditFlashcard(ctx context.Context, args *rest.EditFlashcardRequest) error {
	return nil
}
