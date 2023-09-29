package entities

import (
	"encoding/json"

	"github.com/google/uuid"
)

type (
	User struct {
		Id       uuid.UUID `json:id`
		Login    string    `json:login`
		Password string    `json:password`
	}

	Flashcard struct {
		NativeLanguage string   `json:"native_lang"`
		TargetLang     string   `json:"target_lang"`
		Meaning        string   `json:"word_in_native"`
		Word           string   `json:"word_in_target"`
		UsageExamples  []string `json:"usage"`
	}

	Deck struct {
		Id    string    `json:id`
		Name  string    `json:name`
		Owner uuid.UUID `json:owner`
	}
)

func (u *User) ToJson() ([]byte, error) {
	return json.Marshal(u)
}

func (u *Flashcard) ToJson() ([]byte, error) {
	return json.Marshal(u)
}

func (u *Deck) ToJson() ([]byte, error) {
	return json.Marshal(u)
}
