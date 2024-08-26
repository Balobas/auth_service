-- +goose Up

create type user_role as enum('user', 'admin');

create table users (
    uid varchar not null primary key,
    name varchar,
    email varchar,
    phone varchar,
    password varchar,
    role user_role,
    created_at timestamp,
    updated_at timestamp
);

create type user_permission as enum('not_verified', 'base')

create table user_permissions (
    user_uid varchar primary key,
    permissions []user_permission,
    foreign key(user_uid) references users(uid)
);

create table user_devices (
    uid varchar not null primary key,
    user_uid varchar,
    name varchar,
    os varchar,
    connected_at varchar,
    foreign key(user_uid) references users(uid)
);

create table sessions (
    uid varchar not null primary key,
    user_uid varchar not null,
    device_uid varchar not null,
    created_at timestamp,
    updated_at timestamp,
    foreign key(user_uid) references users(uid),
    foreign key(device_uid) references user_devices(uid)
);

create table black_list_users (
    user_uid varchar not null primary key,
    reason varchar,
    foreign key(user_uid) references users(uid)
);

create table black_list_devices (
    device_uid varchar not null primary key,
    reason varchar,
    foreign key(device_uid) references user_devices(uid)
);

create table verification (
    user_uid varchar not null primary key,
    token varchar not null,
    foreign key(user_uid) references users(uid)
);

-- +goose Down

drop table verification;
drop table black_list_devices;
drop table black_list_users;
drop table sessions;
drop table user_devices;
drop table user_permissions;
drop type user_permission;
drop table users;
drop type user_role;
