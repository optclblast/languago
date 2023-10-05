package entities

import (
	"encoding/json"
	// "languago/pkg/models/requests/rest"
	// "languago/pkg/repository/postgresql"

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

func (m *User) ToJson() ([]byte, error) {
	return json.Marshal(m)
}

// func (m *User) ToModel(v any) error {
// 	switch vType := v.(type) {
// 	case postgresql.User:
// 	// case mysql.User:
// 	case rest.CreateUserRequest:
// 	}
// }

func (m *Flashcard) ToJson() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Deck) ToJson() ([]byte, error) {
	return json.Marshal(m)
}
