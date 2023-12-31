package rest

import (
	"languago/pkg/models/entities"
)

type (
	NewFlashcardRequest struct {
		NativeLanguage string `json:"native_lang"`
		TargetLang     string `json:"target_lang"`
		Content        struct {
			WordInNative  string   `json:"word_in_native"`
			WordInTarget  string   `json:"word_in_target"`
			UsageExamples []string `json:"usage"`
		} `json:"content"`
	}

	NewFlashcardResponse struct {
		Errors []string `json:"errors,omitempty"` // May be empty in OK
	}

	GetFlashcardResponse struct {
		Flashcards []*entities.Flashcard `json:"flashcards"`
	}

	EditFlashcardRequest struct {
		Id            string   `json:"id"`
		WordInNative  string   `json:"word_in_native,omitempty"`
		WordInTarget  string   `json:"word_in_target,omitempty"`
		UsageExamples []string `json:"usage,omitempty"`
	}
)

// TODO grammar cards
