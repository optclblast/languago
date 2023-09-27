-- User 
-- name: CreateUser :exec
INSERT INTO users 
    (id, login, password) 
    VALUES 
    ($1, $2, $3);

-- name: SelectUser :one
SELECT * FROM users 
    WHERE id = $1 AND login = $2;

-- name: SelectUserByLogin :one
SELECT * FROM users 
    WHERE login = $1;

-- name: SelectUserByID :one
SELECT * FROM users 
    WHERE id = $1;

-- name: UpdateUserLogin :one
UPDATE users SET login = $1
    WHERE id = $2
    RETURNING id, login;

-- name: UpdateUserPassword :exec
UPDATE users SET password = $1
    WHERE id = $2;

-- name: DeleteUser :exec 
DELETE FROM users 
    WHERE id = $1;

-- Flashcards
-- name: CreateFlashcard :one
INSERT INTO flashcards
    (id, word, meaning, usage)
    VALUES
    ($1, $2, $3, $4)
    RETURNING word, meaning, usage;

-- name: SelectFlashcardByID :one
SELECT * FROM flashcards 
    WHERE id = $1;

-- name: SelectFlashcardByMeaning :many
SELECT * FROM flashcards AS f
    INNER JOIN flashcard_decks AS d
        ON d.deck_id = $1
    WHERE meaning = $2;

-- name: SelectFlashcardByWord :many
SELECT * FROM flashcards AS f
    INNER JOIN flashcard_decks AS d
        ON d.deck_id = $1
    WHERE word = $2;

-- name: UpdateFlashcard :exec
UPDATE flashcards SET
    word = $1,
    meaning = $2,
    usage = $3
    WHERE id = $4;

-- name: DeleteFlashcard :exec
DELETE FROM flashcards 
    WHERE id = $1;

-- Decks
-- name: CreateDeck :one
INSERT INTO decks 
    (id, name, owner)
    VALUES
    ($1, $2, $3)
    RETURNING name, owner;

-- name: SelectOwnerDecks :many
SELECT * FROM decks
    WHERE owner = $1;

-- name: SelectDeck :one
SELECT * FROM decks 
    WHERE id = $1;

-- name: SelectDecksByName :many
SELECT * FROM decks
    WHERE name = $1;

-- name: EditDeckProps :exec
UPDATE decks SET
    name = $1
    WHERE id = $2;

-- name: DeleteDeck :exec
DELETE FROM decks
    WHERE id = $1;

-- name: AddToDeck :exec
INSERT INTO flashcard_decks
    (deck_id, flashcard_id)
    VALUES
    ($1, $2); 

-- name: DeleteFromDeck :exec
DELETE FROM flashcard_decks
    WHERE flashcard_id = $1 AND
        deck_id = $2;
    