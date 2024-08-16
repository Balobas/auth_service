-- +goose Up

create type user_role as enum('user', 'admin');

create table users (
    id serial primary key,
    name varchar,
    email varchar,
    password varchar,
    password_confirm varchar,
    role user_role,
    created_at timestamp,
    updated_at timestamp
);

-- +goose Down
drop table users;
drop type user_role;
