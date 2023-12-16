drop schema if exists languago cascade;
create schema languago;

drop table if exists languago.users;
create table languago.users
(
    id       uuid not null primary key,
    login    varchar(100),
    password text,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    deleted_at timestamp default null
);

drop index if exists index_users_login;
create index index_users_login
    on languago.users (login);

drop index if exists index_users_id_login;
create index index_users_id_login
    on languago.users (id, login);

drop table if exists languago.flashcards;
create table languago.flashcards
(
    id      uuid not null primary key,
    word    text,
    meaning text,
    usage   text[]
);
drop index if exists index_flashcards_word_meaning;
drop index if exists index_flashcards_word;
drop index if exists index_flashcards_meaning;

create index index_flashcards_word
    on languago.flashcards (word);
create index index_flashcards_meaning
    on languago.flashcards (meaning);
create index index_flashcards_word_meaning
    on languago.flashcards (word, meaning);


drop table if exists languago.decks;
create table languago.decks
(
    id    uuid not null primary key,
    name  varchar(200),
    owner uuid references languago.users(id)
);

drop index if exists index_decks_decks_name_owner;
drop index if exists index_decks_decks_name;
drop index if exists index_decks_decks_owner;

create index index_decks_decks_name_owner
    on languago.decks (name, owner);
create index index_decks_decks_name
    on languago.decks (name);
create index index_decks_decks_owner
    on languago.decks (owner);


drop table if exists languago.flashcard_decks;
create table languago.flashcard_decks
(
    deck_id      uuid references languago.decks(id),
    flashcard_id uuid references languago.flashcards(id)
);
drop index if exists index_flashcard_decks_deck_id;
drop index if exists index_flashcard_decks_flashcard_id;
drop index if exists index_flashcard_decks_deck_id_flashcard_id;

create index index_flashcard_decks_deck_id
    on languago.flashcard_decks (deck_id);
create index index_flashcard_decks_flashcard_id
    on languago.flashcard_decks (flashcard_id);
create index index_flashcard_decks_deck_id_flashcard_id
    on languago.flashcard_decks (deck_id, flashcard_id);

drop table if exists languago.profiles;
create table languago.profiles (
                                   id uuid not null primary key,
                                   name varchar(100),
                                   description varchar(550),
                                   avatar_source text,
                                   user_id uuid references languago.users(id)
);
create index index_profiles_name on languago.profiles
    (name);
create index index_profiles_user_id on languago.profiles
    (user_id);
create index index_profiles_user_id_name on languago.profiles
    (user_id, name);


drop table if exists languago.classes;
create table languago.classes (
                                  id uuid primary key,
                                  name varchar(150),
                                  description varchar(550)
);

alter table languago.flashcards add column   created_at timestamp default now();
alter table languago.flashcards add column    updated_at timestamp default now();
alter table languago.flashcards add column   deleted_at timestamp default null;


drop table if exists languago.flashcard_answers_statistics;
create table languago.flashcard_answers_statistics (
                                                       flashcard_id uuid references languago.flashcards(id),
                                                       user_id uuid references languago.users(id),
                                                       deck_id uuid references languago.decks(id),
                                                       last_answer_at timestamp default now(),
                                                       correct_answers_count bigint default 0,
                                                       wrong_answers_count bigint default 0,

                                                       primary key(flashcard_id, user_id, deck_id)
);

drop table if exists languago.flashcard_answers_history;
create table languago.flashcard_answers_history(
                                                   id uuid primary key,
                                                   flashcard_id uuid references languago.flashcards(id),
                                                   user_id uuid references languago.users(id),
                                                   deck_id uuid references languago.decks(id),
                                                   correct boolean default false,
                                                   answered_at timestamp default now()
);
create index index_answers_history_flashcard_id_user_id_deck_id on
    flashcard_answers_history (flashcard_id, user_id, deck_id);

drop table if exists languago.sessions;
create table languago.sessions (
    id uuid not null,
    user_id uuid not null,
    created_at timestamp default now(),
    expired_at timestamp not null,
    token text,
    primary key (id, user_id)
);

drop index if exists index_sessions_token;
create index index_sessions_token
    on languago.sessions (token);
