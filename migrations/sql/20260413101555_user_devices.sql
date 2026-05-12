-- +goose Up
-- +goose StatementBegin

create type device_type as enum('user_device', 'system_device');

create table devices (
    uid uuid not null,
    maintainer_uid uuid not null,
    type device_type not null,
    name varchar,
    agent varchar,
    language varchar,
    created_at timestamp not null default now(),
    authorized_at timestamp not null default now(),
    unauthorized_at timestamp,
    primary key(uid, maintainer_uid)
);

alter table sessions add column device_uid uuid not null;
-- ALTER TABLE sessions 
-- ADD CONSTRAINT device_uid_fk 
-- FOREIGN KEY (device_uid) 
-- REFERENCES users_devices (uid);
-- TODO: add constraint

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table sessions drop column device_uid;
drop table users_devices;
drop type device_type;
-- +goose StatementEnd
