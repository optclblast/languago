package entities

import (
	"encoding/json"
	"languago/infrastructure/repository/postgresql"
	"languago/pkg/models"

	// "languago/models/requests/rest"
	// "languago/repository/postgresql"

	"github.com/google/uuid"
)

type (
	User struct {
		Id       uuid.UUID `json:"id"`
		Login    string    `json:"login"`
		Password string    `json:"password"`
	}

	Flashcard struct {
		ID             uuid.UUID `json:"id"`
		NativeLanguage string    `json:"native_lang"`
		TargetLang     string    `json:"target_lang"`
		Meaning        string    `json:"word_in_native"`
		Word           string    `json:"word_in_target"`
		UsageExamples  []string  `json:"usage"`
	}

	Deck struct {
		Id    string    `json:"id"`
		Name  string    `json:"name"`
		Owner uuid.UUID `json:"owner"`
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

func UserFromPG(user postgresql.User) *User {
	return &User{
		Id:       user.ID,
		Login:    user.Login.String,
		Password: user.Password.String,
	}
}

func (u *User) ToModel() *models.User {
	return &models.User{
		Id:    u.Id,
		Login: u.Login,
	}
}
