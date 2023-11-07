package flashcards

import (
	"context"
	"languago/internal/pkg/logger"
	"languago/internal/pkg/models/requests/rest"
	"languago/internal/pkg/repository"

	"github.com/google/uuid"
)

type FlashcardsController interface {
	CreateFlashcard(ctx context.Context, req *rest.NewFlashcardRequest) (*rest.NewFlashcardResponse, error)
	GetFlashcard(ctx context.Context, args GetFlashcardParams) (*rest.GetFlashcardResponse, error)
	DeleteFlashcard()
	EditFlashcard()
}

type flashcardController struct {
	log     logger.Logger
	storage repository.DatabaseInteractor
}

func NewFlashcardsController(
	log logger.Logger,
	storage repository.DatabaseInteractor,
) FlashcardsController {
	return &flashcardController{
		log:     log,
		storage: storage,
	}
}

func (c *flashcardController) CreateFlashcard() {}

type GetFlashcardParams struct {
	Id      uuid.UUID
	DeckId  uuid.UUID
	Word    string
	Meaning string
}

// func (c *flashcardController) GetFlashcard()    {}
// func (c *flashcardController) DeleteFlashcard() {}
// func (c *flashcardController) EditFlashcard()   {}
