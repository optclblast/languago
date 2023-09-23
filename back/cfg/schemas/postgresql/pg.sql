CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "login" varchar(100),
  "password" text
);

CREATE TABLE "flashcards" (
  "id" uuid PRIMARY KEY,
  "word" text,
  "meaning" text,
  "usage" text[]
);

CREATE TABLE "flashcard_decks" (
  "deck_id" uuid,
  "flashcard_id" uuid
);

CREATE TABLE "decks" (
  "id" uuid PRIMARY KEY,
  "name" varchar(200),
  "owner" uuid
);

ALTER TABLE "decks" ADD FOREIGN KEY ("id") REFERENCES "flashcard_decks" ("deck_id");

ALTER TABLE "flashcards" ADD FOREIGN KEY ("id") REFERENCES "flashcard_decks" ("flashcard_id");

ALTER TABLE "decks" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");
