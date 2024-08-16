-- +goose Up
Create table token (
    id int PRIMARY KEY,
    value varchar,
    name varchar,
    expired_at timestamp,
    created_at timestamp
);

-- +goose Down
DROP table token;

