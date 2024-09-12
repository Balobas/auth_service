-- +goose Up

create type user_role as enum('user', 'admin');

create table users (
    uid varchar not null primary key,
    email varchar,
    role user_role,
    created_at timestamp,
    updated_at timestamp
);

create table users_credentials (
    user_uid not null primary key,
    h_password varchar,
    foreign key(user_uid) references users(uid)
);

create type user_permission as enum('not_verified', 'base');

create table user_permissions (
    user_uid varchar primary key,
    permissions []user_permission,
    foreign key(user_uid) references users(uid)
);

create table sessions (
    uid varchar not null primary key,
    user_uid varchar not null,
    created_at timestamp,
    updated_at timestamp,
    foreign key(user_uid) references users(uid)
);

create type verification_status as enum('created', 'waiting');

create table verification (
    user_uid varchar not null primary key,
    token varchar not null,
    status verification_status not null,
    created_at timestamp,
    foreign key(user_uid) references users(uid)
);

-- +goose Down

drop table verification;
drop type verification_status;
drop table sessions;
drop table user_devices;
drop table user_permissions;
drop type user_permission;
drop table users_credentials;
drop table users;
drop type user_role;
