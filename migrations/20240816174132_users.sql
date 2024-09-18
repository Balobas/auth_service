-- +goose Up

create type user_role as enum('user', 'admin');

create table users (
    uid uuid not null primary key,
    email varchar,
    role user_role,
    created_at timestamp,
    updated_at timestamp
);

create table users_credentials (
    user_uid uuid not null primary key,
    h_password varchar,
    foreign key(user_uid) references users(uid) on delete cascade
);

create type user_permission as enum('not_verified', 'base');

create table user_permissions (
    user_uid uuid primary key,
    permissions user_permission[],
    foreign key(user_uid) references users(uid) on delete cascade
);

create table sessions (
    uid uuid not null primary key,
    user_uid uuid not null,
    created_at timestamp,
    updated_at timestamp,
    foreign key(user_uid) references users(uid) on delete cascade
);

create type verification_status as enum('created', 'waiting');

create table verification (
    user_uid uuid not null primary key,
    email varchar not null unique,
    token varchar not null,
    status verification_status not null,
    created_at timestamp,
    updated_at timestamp,
    foreign key(user_uid) references users(uid) on delete cascade
);

create table config (
    key varchar not null,
    value json
);

insert into config(key, value) values
('min_password_len', '6'),
('access_jwt_ttl', '"1h"'),
('refresh_jwt_ttl', '"24h"'),
('verification_token_len', '16'),
('send_verification_interval', '"3m"'),
('verification_worker_batch_size', '10'),
('email_verification_template', '"Подтвердите вашу почту перейдя по ссылке {{.Scheme}}\/{{.Token }}"'),
('http_verification_scheme', 'null');

-- +goose Down

drop table config;
drop table verification;
drop type verification_status;
drop table sessions;
drop table user_permissions;
drop type user_permission;
drop table users_credentials;
drop table users;
drop type user_role;
