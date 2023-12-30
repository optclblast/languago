\c '\set AUTOCOMMIT on'
drop database if exists FOODB;
CREATE DATABASE FOODB;

\c FOODB;

CREATE SCHEMA BAR;

SET search_path TO BAR;

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(255) NOT NULL,
                       email VARCHAR(255) NOT NULL
);

CREATE INDEX index_users_id ON BAR.users (id);