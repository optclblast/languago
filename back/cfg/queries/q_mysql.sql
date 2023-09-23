-- User 
-- name: CreateUser :exec
INSERT INTO users 
    (id, login, password) 
    VALUES 
    ($1, $2, $3);

-- name: SelectUser :exec
SELECT * FROM users 
    WHERE id = $1 AND login = $2;

-- name: UpdateUserLogin :exec
UPDATE users SET login = $1
    WHERE id = $2;

-- name: UpdateUserPassword :exec
UPDATE users SET password = $1
    WHERE id = $2;

-- name: DeleteUser :exec 
DELETE FROM users 
    WHERE id = $1;

-- Flashcards
-- name: CreateFlashcard :exec
INSERT INTO flashcards 
    (id, word, meaning, usage) 
    VALUES ($1, $2, $3, $4);

-- name: SelectFlashcards :many
SELECT * FROM flashcards 
    WHERE id = $1 AND 
    word = $2 AND 
    meaning = $3 AND 
    usage = $4;

-- name: UpdateFlashcard :exec
UPDATE flashcards SET
    id = $1,
    word = $2,
    meaning = $3,
    usage = $4
    WHERE id = $5;

-- name: DeleteFlashcard :exec
DELETE FROM flashcards 
    WHERE id = $1;

-- Decks
-- name: CreateDeck :exec
INSERT INTO decks 
    (id, name, owner)
    VALUES
    ($1, $2, $3);

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
    