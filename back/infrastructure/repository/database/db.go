// TODO

package database

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type DBTX interface {
	sq.Runner
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type databaseController struct {
	conn *sql.DB
}

type AddToDeckParams struct {
	DeckID      uuid.UUID
	FlashcardID uuid.UUID
}

func (c *databaseController) AddToDeck(ctx context.Context, arg AddToDeckParams) error {
	stmt := sq.Insert("flashcard_decks").Columns(
		"deck_id", "flashcard_id",
	).Values(
		arg.DeckID, arg.FlashcardID,
	)

	_, err := stmt.RunWith(c.conn).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("error add flashcard to deck: %w", err)
	}

	return nil
}

type CreateDeckParams struct {
	ID    uuid.UUID
	Name  string
	Owner uuid.UUID
}

// Decks
func (c *databaseController) CreateDeck(ctx context.Context, arg CreateDeckParams) error {
	_, err := sq.Insert("decks").Columns(
		"id",
		"name",
		"owner",
	).Values(
		arg.ID,
		arg.Name,
		arg.Owner,
	).RunWith(c.conn).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("error create deck: %w", err)
	}

	return nil
}

type CreateFlashcardParams struct {
	ID      uuid.UUID
	Word    string
	Meaning string
	Usage   []string
}

// Flashcards
func (c *databaseController) CreateFlashcard(ctx context.Context, arg CreateFlashcardParams) error {
	_, err := sq.Insert("flashcards").Columns(
		"id",
		"word",
		"meaning",
		"usage",
	).Values(
		arg.ID,
		arg.Word,
		arg.Meaning,
		arg.Usage,
	).RunWith(c.conn).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("error create flashcard: %w", err)
	}

	return nil
}

type CreateUserParams struct {
	ID       uuid.UUID
	Login    string
	Password string
}

func (c *databaseController) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := sq.Insert("users").Columns(
		"id",
		"login",
		"password",
	).Values(
		arg.ID,
		arg.Login,
		arg.Password,
	).RunWith(c.conn).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("error create user: %w", err)
	}

	return nil
}

func (c *databaseController) DeleteDeck(ctx context.Context, id uuid.UUID) error {
	_, err := sq.Delete("decks").Where(sq.Eq{"id": id}).RunWith(c.conn).ExecContext(ctx)

	return err
}

func (c *databaseController) DeleteFlashcard(ctx context.Context, id uuid.UUID) error {
	_, err := sq.Delete("flashcards").Where(sq.Eq{"id": id}).RunWith(c.conn).ExecContext(ctx)

	return err
}

type DeleteFromDeckParams struct {
	FlashcardID uuid.UUID
	DeckID      uuid.UUID
}

func (c *databaseController) DeleteFromDeck(ctx context.Context, arg DeleteFromDeckParams) error {
	_, err := sq.Delete("flashcard_decks").Where(
		sq.Eq{"deck_id": arg.DeckID},
		sq.Eq{"flashcard_id": arg.FlashcardID},
	).RunWith(c.conn).ExecContext(ctx)

	return err
}

func (c *databaseController) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := sq.Delete("users").Where(
		sq.Eq{"id": id},
	).RunWith(c.conn).ExecContext(ctx)

	return err
}

// const editDeckProps = `-- name: EditDeckProps :exec
// UPDATE decks SET
//     name = $1
//     WHERE id = $2
// `

type EditDeckPropsParams struct {
	Name string
	ID   uuid.UUID
}

func (c *databaseController) EditDeckProps(ctx context.Context, arg EditDeckPropsParams) error {
	_, err := sq.Update("decks").SetMap(
		sq.Eq{
			"name": arg.Name,
		},
	).Where(sq.Eq{"id": arg.ID}).RunWith(c.conn).ExecContext(ctx)

	return err
}

// const selectDeck = `-- name: SelectDeck :one
// SELECT id, name, owner FROM decks
//     WHERE id = $1
// `

// func (c *databaseController) SelectDeck(ctx context.Context, id uuid.UUID) (Deck, error) {
// 	row := c.conn.QueryRowContext(ctx, selectDeck, id)
// 	var i Deck
// 	err := row.Scan(&i.ID, &i.Name, &i.Owner)
// 	return i, err
// }

// const selectDecksByName = `-- name: SelectDecksByName :many
// SELECT id, name, owner FROM decks
//     WHERE name = $1
// `

// func (c *databaseController) SelectDecksByName(ctx context.Context, name sql.NullString) ([]Deck, error) {
// 	rows, err := c.conn.QueryContext(ctx, selectDecksByName, name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var items []Deck
// 	for rows.Next() {
// 		var i Deck
// 		if err := rows.Scan(&i.ID, &i.Name, &i.Owner); err != nil {
// 			return nil, err
// 		}
// 		items = append(items, i)
// 	}
// 	if err := rows.Close(); err != nil {
// 		return nil, err
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return items, nil
// }

// const selectFlashcardByID = `-- name: SelectFlashcardByID :one
// SELECT id, word, meaning, usage FROM flashcards
//     WHERE id = $1
// `

// func (c *databaseController) SelectFlashcardByID(ctx context.Context, id uuid.UUID) (Flashcard, error) {
// 	row := c.conn.QueryRowContext(ctx, selectFlashcardByID, id)
// 	var i Flashcard
// 	err := row.Scan(
// 		&i.ID,
// 		&i.Word,
// 		&i.Meaning,
// 		pq.Array(&i.Usage),
// 	)
// 	return i, err
// }

// const selectFlashcardByMeaning = `-- name: SelectFlashcardByMeaning :many
// SELECT id, word, meaning, usage, deck_id, flashcard_id FROM flashcards AS f
//     INNER JOIN flashcard_decks AS d
//         ON d.deck_id = $1
//     WHERE meaning = $2
// `

// type SelectFlashcardByMeaningParams struct {
// 	DeckID  uuid.NullUUID  `db:"deck_id" json:"deck_id"`
// 	Meaning sql.NullString `db:"meaning" json:"meaning"`
// }

// type SelectFlashcardByMeaningRow struct {
// 	ID          uuid.UUID      `db:"id" json:"id"`
// 	Word        sql.NullString `db:"word" json:"word"`
// 	Meaning     sql.NullString `db:"meaning" json:"meaning"`
// 	Usage       []string       `db:"usage" json:"usage"`
// 	DeckID      uuid.NullUUID  `db:"deck_id" json:"deck_id"`
// 	FlashcardID uuid.NullUUID  `db:"flashcard_id" json:"flashcard_id"`
// }

// func (c *databaseController) SelectFlashcardByMeaning(ctx context.Context, arg SelectFlashcardByMeaningParams) ([]SelectFlashcardByMeaningRow, error) {
// 	rows, err := c.conn.QueryContext(ctx, selectFlashcardByMeaning, arg.DeckID, arg.Meaning)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var items []SelectFlashcardByMeaningRow
// 	for rows.Next() {
// 		var i SelectFlashcardByMeaningRow
// 		if err := rows.Scan(
// 			&i.ID,
// 			&i.Word,
// 			&i.Meaning,
// 			pq.Array(&i.Usage),
// 			&i.DeckID,
// 			&i.FlashcardID,
// 		); err != nil {
// 			return nil, err
// 		}
// 		items = append(items, i)
// 	}
// 	if err := rows.Close(); err != nil {
// 		return nil, err
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return items, nil
// }

// const selectFlashcardByWord = `-- name: SelectFlashcardByWord :many
// SELECT id, word, meaning, usage, deck_id, flashcard_id FROM flashcards AS f
//     INNER JOIN flashcard_decks AS d
//         ON d.deck_id = $1
//     WHERE word = $2
// `

// type SelectFlashcardByWordParams struct {
// 	DeckID uuid.NullUUID  `db:"deck_id" json:"deck_id"`
// 	Word   sql.NullString `db:"word" json:"word"`
// }

// type SelectFlashcardByWordRow struct {
// 	ID          uuid.UUID      `db:"id" json:"id"`
// 	Word        sql.NullString `db:"word" json:"word"`
// 	Meaning     sql.NullString `db:"meaning" json:"meaning"`
// 	Usage       []string       `db:"usage" json:"usage"`
// 	DeckID      uuid.NullUUID  `db:"deck_id" json:"deck_id"`
// 	FlashcardID uuid.NullUUID  `db:"flashcard_id" json:"flashcard_id"`
// }

// func (c *databaseController) SelectFlashcardByWord(ctx context.Context, arg SelectFlashcardByWordParams) ([]SelectFlashcardByWordRow, error) {
// 	rows, err := c.conn.QueryContext(ctx, selectFlashcardByWord, arg.DeckID, arg.Word)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var items []SelectFlashcardByWordRow
// 	for rows.Next() {
// 		var i SelectFlashcardByWordRow
// 		if err := rows.Scan(
// 			&i.ID,
// 			&i.Word,
// 			&i.Meaning,
// 			pq.Array(&i.Usage),
// 			&i.DeckID,
// 			&i.FlashcardID,
// 		); err != nil {
// 			return nil, err
// 		}
// 		items = append(items, i)
// 	}
// 	if err := rows.Close(); err != nil {
// 		return nil, err
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return items, nil
// }

// const selectOwnerDecks = `-- name: SelectOwnerDecks :many
// SELECT id, name, owner FROM decks
//     WHERE owner = $1
// `

// func (c *databaseController) SelectOwnerDecks(ctx context.Context, owner uuid.NullUUID) ([]Deck, error) {
// 	rows, err := c.conn.QueryContext(ctx, selectOwnerDecks, owner)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var items []Deck
// 	for rows.Next() {
// 		var i Deck
// 		if err := rows.Scan(&i.ID, &i.Name, &i.Owner); err != nil {
// 			return nil, err
// 		}
// 		items = append(items, i)
// 	}
// 	if err := rows.Close(); err != nil {
// 		return nil, err
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return items, nil
// }

// const selectUser = `-- name: SelectUser :one
// SELECT id, login, password FROM users
//     WHERE id = $1 AND login = $2
// `

// type SelectUserParams struct {
// 	ID    uuid.UUID      `db:"id" json:"id"`
// 	Login sql.NullString `db:"login" json:"login"`
// }

// func (c *databaseController) SelectUser(ctx context.Context, arg SelectUserParams) (User, error) {
// 	row := c.conn.QueryRowContext(ctx, selectUser, arg.ID, arg.Login)
// 	var i User
// 	err := row.Scan(&i.ID, &i.Login, &i.Password)
// 	return i, err
// }

// const selectUserByID = `-- name: SelectUserByID :one
// SELECT id, login, password FROM users
//     WHERE id = $1
// `

// func (c *databaseController) SelectUserByID(ctx context.Context, id uuid.UUID) (User, error) {
// 	row := c.conn.QueryRowContext(ctx, selectUserByID, id)
// 	var i User
// 	err := row.Scan(&i.ID, &i.Login, &i.Password)
// 	return i, err
// }

// const selectUserByLogin = `-- name: SelectUserByLogin :one
// SELECT id, login, password FROM users
//     WHERE login = $1
// `

// func (c *databaseController) SelectUserByLogin(ctx context.Context, login sql.NullString) (User, error) {
// 	row := c.conn.QueryRowContext(ctx, selectUserByLogin, login)
// 	var i User
// 	err := row.Scan(&i.ID, &i.Login, &i.Password)
// 	return i, err
// }

// const updateFlashcard = `-- name: UpdateFlashcard :exec
// UPDATE flashcards SET
//     word = $1,
//     meaning = $2,
//     usage = $3
//     WHERE id = $4
// `

// type UpdateFlashcardParams struct {
// 	Word    sql.NullString `db:"word" json:"word"`
// 	Meaning sql.NullString `db:"meaning" json:"meaning"`
// 	Usage   []string       `db:"usage" json:"usage"`
// 	ID      uuid.UUID      `db:"id" json:"id"`
// }

// func (c *databaseController) UpdateFlashcard(ctx context.Context, arg UpdateFlashcardParams) error {
// 	_, err := c.conn.ExecContext(ctx, updateFlashcard,
// 		arg.Word,
// 		arg.Meaning,
// 		pq.Array(arg.Usage),
// 		arg.ID,
// 	)
// 	return err
// }

// const updateUserLogin = `-- name: UpdateUserLogin :one
// UPDATE users SET login = $1
//     WHERE id = $2
//     RETURNING id, login
// `

// type UpdateUserLoginParams struct {
// 	Login sql.NullString `db:"login" json:"login"`
// 	ID    uuid.UUID      `db:"id" json:"id"`
// }

// type UpdateUserLoginRow struct {
// 	ID    uuid.UUID      `db:"id" json:"id"`
// 	Login sql.NullString `db:"login" json:"login"`
// }

// func (c *databaseController) UpdateUserLogin(ctx context.Context, arg UpdateUserLoginParams) (UpdateUserLoginRow, error) {
// 	row := c.conn.QueryRowContext(ctx, updateUserLogin, arg.Login, arg.ID)
// 	var i UpdateUserLoginRow
// 	err := row.Scan(&i.ID, &i.Login)
// 	return i, err
// }

// const updateUserPassword = `-- name: UpdateUserPassword :exec
// UPDATE users SET password = $1
//     WHERE id = $2
// `

// type UpdateUserPasswordParams struct {
// 	Password sql.NullString `db:"password" json:"password"`
// 	ID       uuid.UUID      `db:"id" json:"id"`
// }

// func (c *databaseController) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
// 	_, err := c.conn.ExecContext(ctx, updateUserPassword, arg.Password, arg.ID)
// 	return err
// }
