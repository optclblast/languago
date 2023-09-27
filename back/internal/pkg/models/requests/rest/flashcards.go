package rest

import (
	"fmt"
	"languago/internal/pkg/repository/postgresql"
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
		Flashcards []*Flashcard `json:"flashcards"`
	}

	Flashcard struct {
		NativeLanguage string   `json:"native_lang"`
		TargetLang     string   `json:"target_lang"`
		WordInNative   string   `json:"word_in_native"`
		WordInTarget   string   `json:"word_in_target"`
		UsageExamples  []string `json:"usage"`
	}

	EditFlashcardRequest struct {
		Id            string   `json:"id"`
		WordInNative  string   `json:"word_in_native,omitempty"`
		WordInTarget  string   `json:"word_in_target,omitempty"`
		UsageExamples []string `json:"usage,omitempty"`
	}
)

// TODO grammar cards

// utils

func (r *GetFlashcardResponse) FromFlashcardObject(f []*postgresql.Flashcard) error {
	if f == nil {
		return fmt.Errorf("error object is empty.")
	}
	var obj []*Flashcard = make([]*Flashcard, 0, len(f))
	for _, cardRaw := range f {
		var fc *Flashcard
		fc.WordInNative = cardRaw.Meaning.String
		fc.WordInTarget = cardRaw.Word.String
		fc.UsageExamples = cardRaw.Usage
		obj = append(obj, fc)
	}
	r.Flashcards = obj
	return nil
}

func (r *GetFlashcardResponse) FromFlashcardByWordObject(f []postgresql.SelectFlashcardByWordRow) error {
	if f == nil {
		return fmt.Errorf("error object is empty.")
	}
	var obj []*Flashcard = make([]*Flashcard, 0, len(f))
	for _, cardRaw := range f {
		var fc *Flashcard
		fc.WordInNative = cardRaw.Meaning.String
		fc.WordInTarget = cardRaw.Word.String
		fc.UsageExamples = cardRaw.Usage
		obj = append(obj, fc)
	}
	r.Flashcards = obj
	return nil
}

func (r *GetFlashcardResponse) FromFlashcardByMeaningObject(f []postgresql.SelectFlashcardByMeaningRow) error {
	if f == nil {
		return fmt.Errorf("error object is empty.")
	}
	var obj []*Flashcard = make([]*Flashcard, 0, len(f))
	for _, cardRaw := range f {
		var fc *Flashcard
		fc.WordInNative = cardRaw.Meaning.String
		fc.WordInTarget = cardRaw.Word.String
		fc.UsageExamples = cardRaw.Usage
		obj = append(obj, fc)
	}
	r.Flashcards = obj
	return nil
}
