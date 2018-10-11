
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE DATABASE ivy_winter

CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(80) NOT NULL,
  password VARCHAR(80) NOT NULL,
  token VARCHAR(80) NOT NULL
);

CREATE TABLE IF NOT EXISTS photos (
  id SERIAL NOT NULL,
  data VARCHAR(255) NOT NULL
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE users, photos;
