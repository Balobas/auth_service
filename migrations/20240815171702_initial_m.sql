-- +goose Up

CREATE TABLE auth (
    id int PRIMARY KEY,
    token varchar
);

-- +goose Down
DROP TABLE auth;
