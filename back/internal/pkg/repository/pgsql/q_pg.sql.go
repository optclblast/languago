package pgsql

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const addToDeck = `-- name: AddToDeck :exec
INSERT INTO flashcard_decks
    (deck_id, flashcard_id)
    VALUES
    ($1, $2)
`

type AddToDeckParams struct {
	DeckID      uuid.NullUUID `db:"deck_id" json:"deck_id"`
	FlashcardID uuid.NullUUID `db:"flashcard_id" json:"flashcard_id"`
}

func (q *Queries) AddToDeck(ctx context.Context, arg AddToDeckParams) error {
	_, err := q.db.ExecContext(ctx, addToDeck, arg.DeckID, arg.FlashcardID)
	return err
}

const createDeck = `-- name: CreateDeck :one
INSERT INTO decks 
    (id, name, owner)
    VALUES
    ($1, $2, $3)
    RETURNING name, owner
`

type CreateDeckParams struct {
	ID    uuid.UUID      `db:"id" json:"id"`
	Name  sql.NullString `db:"name" json:"name"`
	Owner uuid.NullUUID  `db:"owner" json:"owner"`
}

type CreateDeckRow struct {
	Name  sql.NullString `db:"name" json:"name"`
	Owner uuid.NullUUID  `db:"owner" json:"owner"`
}

// Decks
func (q *Queries) CreateDeck(ctx context.Context, arg CreateDeckParams) (CreateDeckRow, error) {
	row := q.db.QueryRowContext(ctx, createDeck, arg.ID, arg.Name, arg.Owner)
	var i CreateDeckRow
	err := row.Scan(&i.Name, &i.Owner)
	return i, err
}

const createFlashcard = `-- name: CreateFlashcard :one
INSERT INTO flashcards
    (id, word, meaning, usage)
    VALUES
    ($1, $2, $3, $4)
    RETURNING word, meaning, usage
`

type CreateFlashcardParams struct {
	ID      uuid.UUID      `db:"id" json:"id"`
	Word    sql.NullString `db:"word" json:"word"`
	Meaning sql.NullString `db:"meaning" json:"meaning"`
	Usage   []string       `db:"usage" json:"usage"`
}

type CreateFlashcardRow struct {
	Word    sql.NullString `db:"word" json:"word"`
	Meaning sql.NullString `db:"meaning" json:"meaning"`
	Usage   []string       `db:"usage" json:"usage"`
}

// Flashcards
func (q *Queries) CreateFlashcard(ctx context.Context, arg CreateFlashcardParams) (CreateFlashcardRow, error) {
	row := q.db.QueryRowContext(ctx, createFlashcard,
		arg.ID,
		arg.Word,
		arg.Meaning,
		pq.Array(arg.Usage),
	)
	var i CreateFlashcardRow
	err := row.Scan(&i.Word, &i.Meaning, pq.Array(&i.Usage))
	return i, err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users 
    (id, login, password) 
    VALUES 
    ($1, $2, $3)
`

type CreateUserParams struct {
	ID       uuid.UUID      `db:"id" json:"id"`
	Login    sql.NullString `db:"login" json:"login"`
	Password sql.NullString `db:"password" json:"password"`
}

// User
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser, arg.ID, arg.Login, arg.Password)
	return err
}

const deleteDeck = `-- name: DeleteDeck :exec
DELETE FROM decks
    WHERE id = $1
`

func (q *Queries) DeleteDeck(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteDeck, id)
	return err
}

const deleteFlashcard = `-- name: DeleteFlashcard :exec
DELETE FROM flashcards 
    WHERE id = $1
`

func (q *Queries) DeleteFlashcard(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteFlashcard, id)
	return err
}

const deleteFromDeck = `-- name: DeleteFromDeck :exec
DELETE FROM flashcard_decks
    WHERE flashcard_id = $1 AND
        deck_id = $2
`

type DeleteFromDeckParams struct {
	FlashcardID uuid.NullUUID `db:"flashcard_id" json:"flashcard_id"`
	DeckID      uuid.NullUUID `db:"deck_id" json:"deck_id"`
}

func (q *Queries) DeleteFromDeck(ctx context.Context, arg DeleteFromDeckParams) error {
	_, err := q.db.ExecContext(ctx, deleteFromDeck, arg.FlashcardID, arg.DeckID)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users 
    WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const editDeckProps = `-- name: EditDeckProps :exec
UPDATE decks SET
    name = $1
    WHERE id = $2
`

type EditDeckPropsParams struct {
	Name sql.NullString `db:"name" json:"name"`
	ID   uuid.UUID      `db:"id" json:"id"`
}

func (q *Queries) EditDeckProps(ctx context.Context, arg EditDeckPropsParams) error {
	_, err := q.db.ExecContext(ctx, editDeckProps, arg.Name, arg.ID)
	return err
}

const selectDeck = `-- name: SelectDeck :one
SELECT id, name, owner FROM decks 
    WHERE id = $1
`

func (q *Queries) SelectDeck(ctx context.Context, id uuid.UUID) (Deck, error) {
	row := q.db.QueryRowContext(ctx, selectDeck, id)
	var i Deck
	err := row.Scan(&i.ID, &i.Name, &i.Owner)
	return i, err
}

const selectDecksByName = `-- name: SelectDecksByName :many
SELECT id, name, owner FROM decks
    WHERE name = $1
`

func (q *Queries) SelectDecksByName(ctx context.Context, name sql.NullString) ([]Deck, error) {
	rows, err := q.db.QueryContext(ctx, selectDecksByName, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Deck
	for rows.Next() {
		var i Deck
		if err := rows.Scan(&i.ID, &i.Name, &i.Owner); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectFlashcardByID = `-- name: SelectFlashcardByID :one
SELECT id, word, meaning, usage FROM flashcards 
    WHERE id = $1
`

func (q *Queries) SelectFlashcardByID(ctx context.Context, id uuid.UUID) (Flashcard, error) {
	row := q.db.QueryRowContext(ctx, selectFlashcardByID, id)
	var i Flashcard
	err := row.Scan(
		&i.ID,
		&i.Word,
		&i.Meaning,
		pq.Array(&i.Usage),
	)
	return i, err
}

const selectFlashcardByMeaning = `-- name: SelectFlashcardByMeaning :many
SELECT id, word, meaning, usage, deck_id, flashcard_id FROM flashcards AS f
    INNER JOIN flashcard_decks AS d
        ON d.deck_id = $1
    WHERE meaning = $2
`

type SelectFlashcardByMeaningParams struct {
	DeckID  uuid.NullUUID  `db:"deck_id" json:"deck_id"`
	Meaning sql.NullString `db:"meaning" json:"meaning"`
}

type SelectFlashcardByMeaningRow struct {
	ID          uuid.UUID      `db:"id" json:"id"`
	Word        sql.NullString `db:"word" json:"word"`
	Meaning     sql.NullString `db:"meaning" json:"meaning"`
	Usage       []string       `db:"usage" json:"usage"`
	DeckID      uuid.NullUUID  `db:"deck_id" json:"deck_id"`
	FlashcardID uuid.NullUUID  `db:"flashcard_id" json:"flashcard_id"`
}

func (q *Queries) SelectFlashcardByMeaning(ctx context.Context, arg SelectFlashcardByMeaningParams) ([]SelectFlashcardByMeaningRow, error) {
	rows, err := q.db.QueryContext(ctx, selectFlashcardByMeaning, arg.DeckID, arg.Meaning)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SelectFlashcardByMeaningRow
	for rows.Next() {
		var i SelectFlashcardByMeaningRow
		if err := rows.Scan(
			&i.ID,
			&i.Word,
			&i.Meaning,
			pq.Array(&i.Usage),
			&i.DeckID,
			&i.FlashcardID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectFlashcardByWord = `-- name: SelectFlashcardByWord :many
SELECT id, word, meaning, usage, deck_id, flashcard_id FROM flashcards AS f
    INNER JOIN flashcard_decks AS d
        ON d.deck_id = $1
    WHERE word = $2
`

type SelectFlashcardByWordParams struct {
	DeckID uuid.NullUUID  `db:"deck_id" json:"deck_id"`
	Word   sql.NullString `db:"word" json:"word"`
}

type SelectFlashcardByWordRow struct {
	ID          uuid.UUID      `db:"id" json:"id"`
	Word        sql.NullString `db:"word" json:"word"`
	Meaning     sql.NullString `db:"meaning" json:"meaning"`
	Usage       []string       `db:"usage" json:"usage"`
	DeckID      uuid.NullUUID  `db:"deck_id" json:"deck_id"`
	FlashcardID uuid.NullUUID  `db:"flashcard_id" json:"flashcard_id"`
}

func (q *Queries) SelectFlashcardByWord(ctx context.Context, arg SelectFlashcardByWordParams) ([]SelectFlashcardByWordRow, error) {
	rows, err := q.db.QueryContext(ctx, selectFlashcardByWord, arg.DeckID, arg.Word)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SelectFlashcardByWordRow
	for rows.Next() {
		var i SelectFlashcardByWordRow
		if err := rows.Scan(
			&i.ID,
			&i.Word,
			&i.Meaning,
			pq.Array(&i.Usage),
			&i.DeckID,
			&i.FlashcardID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectOwnerDecks = `-- name: SelectOwnerDecks :many
SELECT id, name, owner FROM decks
    WHERE owner = $1
`

func (q *Queries) SelectOwnerDecks(ctx context.Context, owner uuid.NullUUID) ([]Deck, error) {
	rows, err := q.db.QueryContext(ctx, selectOwnerDecks, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Deck
	for rows.Next() {
		var i Deck
		if err := rows.Scan(&i.ID, &i.Name, &i.Owner); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectUser = `-- name: SelectUser :one
SELECT id, login, password FROM users 
    WHERE id = $1 AND login = $2
`

type SelectUserParams struct {
	ID    uuid.UUID      `db:"id" json:"id"`
	Login sql.NullString `db:"login" json:"login"`
}

func (q *Queries) SelectUser(ctx context.Context, arg SelectUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, selectUser, arg.ID, arg.Login)
	var i User
	err := row.Scan(&i.ID, &i.Login, &i.Password)
	return i, err
}

const selectUserByID = `-- name: SelectUserByID :one
SELECT id, login, password FROM users 
    WHERE id = $1
`

func (q *Queries) SelectUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, selectUserByID, id)
	var i User
	err := row.Scan(&i.ID, &i.Login, &i.Password)
	return i, err
}

const selectUserByLogin = `-- name: SelectUserByLogin :one
SELECT id, login, password FROM users 
    WHERE login = $1
`

func (q *Queries) SelectUserByLogin(ctx context.Context, login sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, selectUserByLogin, login)
	var i User
	err := row.Scan(&i.ID, &i.Login, &i.Password)
	return i, err
}

const updateFlashcard = `-- name: UpdateFlashcard :exec
UPDATE flashcards SET
    word = $1,
    meaning = $2,
    usage = $3
    WHERE id = $4
`

type UpdateFlashcardParams struct {
	Word    sql.NullString `db:"word" json:"word"`
	Meaning sql.NullString `db:"meaning" json:"meaning"`
	Usage   []string       `db:"usage" json:"usage"`
	ID      uuid.UUID      `db:"id" json:"id"`
}

func (q *Queries) UpdateFlashcard(ctx context.Context, arg UpdateFlashcardParams) error {
	_, err := q.db.ExecContext(ctx, updateFlashcard,
		arg.Word,
		arg.Meaning,
		pq.Array(arg.Usage),
		arg.ID,
	)
	return err
}

const updateUserLogin = `-- name: UpdateUserLogin :one
UPDATE users SET login = $1
    WHERE id = $2
    RETURNING id, login
`

type UpdateUserLoginParams struct {
	Login sql.NullString `db:"login" json:"login"`
	ID    uuid.UUID      `db:"id" json:"id"`
}

type UpdateUserLoginRow struct {
	ID    uuid.UUID      `db:"id" json:"id"`
	Login sql.NullString `db:"login" json:"login"`
}

func (q *Queries) UpdateUserLogin(ctx context.Context, arg UpdateUserLoginParams) (UpdateUserLoginRow, error) {
	row := q.db.QueryRowContext(ctx, updateUserLogin, arg.Login, arg.ID)
	var i UpdateUserLoginRow
	err := row.Scan(&i.ID, &i.Login)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users SET password = $1
    WHERE id = $2
`

type UpdateUserPasswordParams struct {
	Password sql.NullString `db:"password" json:"password"`
	ID       uuid.UUID      `db:"id" json:"id"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateUserPassword, arg.Password, arg.ID)
	return err
}
