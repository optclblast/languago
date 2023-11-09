package models

import "github.com/google/uuid"

type (
	User struct {
		Id    uuid.UUID
		Login string
	}

	Flashcard struct {
		NativeLanguage string
		TargetLang     string
		Meaning        string
		Word           string
		UsageExamples  []string
	}

	Deck struct {
		Id    string
		Name  string
		Owner uuid.UUID
	}
)
