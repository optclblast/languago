package pgsql

import (
	"database/sql"

	"github.com/google/uuid"
)

type Deck struct {
	ID    uuid.UUID      `db:"id" json:"id"`
	Name  sql.NullString `db:"name" json:"name"`
	Owner uuid.NullUUID  `db:"owner" json:"owner"`
}

type Flashcard struct {
	ID      uuid.UUID      `db:"id" json:"id"`
	Word    sql.NullString `db:"word" json:"word"`
	Meaning sql.NullString `db:"meaning" json:"meaning"`
	Usage   []string       `db:"usage" json:"usage"`
}

type User struct {
	ID       uuid.UUID      `db:"id" json:"id"`
	Login    sql.NullString `db:"login" json:"login"`
	Password sql.NullString `db:"password" json:"password"`
}
