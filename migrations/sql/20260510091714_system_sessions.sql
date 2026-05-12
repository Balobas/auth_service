-- +goose Up
-- +goose StatementBegin
alter table sessions drop column user_uid;
alter table sessions add column maintainer_uid uuid not null;

create type session_type as enum('user_session', 'system_session');
alter table sessions add column type session_type not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
