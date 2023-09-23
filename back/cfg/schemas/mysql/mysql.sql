CREATE TABLE `users` (
  `id` varchar(36) PRIMARY KEY,
  `login` varchar(100),
  `password` text
);

CREATE TABLE `flashcards` (
  `id` varchar(36) PRIMARY KEY,
  `word` text,
  `meaning` text,
  `usage` text
);

CREATE TABLE `flashcard_decks` (
  `deck_id` varchar(36),
  `flashcard_id` varchar(36)
);

CREATE TABLE `decks` (
  `id` varchar(36) PRIMARY KEY,
  `name` varchar(200),
  `owner` varchar(36)
);

ALTER TABLE `decks` ADD FOREIGN KEY (`id`) REFERENCES `flashcard_decks` (`deck_id`);

ALTER TABLE `flashcards` ADD FOREIGN KEY (`id`) REFERENCES `flashcard_decks` (`flashcard_id`);

ALTER TABLE `decks` ADD FOREIGN KEY (`owner`) REFERENCES `users` (`id`);
