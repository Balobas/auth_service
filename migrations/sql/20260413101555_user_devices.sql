-- +goose Up
-- +goose StatementBegin
create table users_devices (
    uid uuid not null,
    user_uid uuid references users(uid),
    name varchar,
    agent varchar,
    language varchar,
    created_at timestamp not null default now(),
    authorized_at timestamp not null default now(),
    unauthorized_at timestamp,
    primary key(uid, user_uid)
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
-- +goose StatementEnd
