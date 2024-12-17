-- +goose Up
-- +goose StatementBegin

create table outbox_messages (
    uid uuid PRIMARY KEY,
    subject_name varchar,
    payload varchar,
    last_error_msg varchar,
    created_at timestamp,
    updated_at timestamp,
    send_at timestamp
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table outbox_messages;

-- +goose StatementEnd
