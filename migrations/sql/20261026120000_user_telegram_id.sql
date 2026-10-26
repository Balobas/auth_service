-- +goose Up
-- +goose StatementBegin

-- Привязка Telegram-аккаунта к пользователю (вход через мини-апп). Связь 1:1:
-- один аккаунт приложения ↔ один Telegram. nullable (у email-юзеров нет).
alter table users add column telegram_id bigint unique;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table users drop column if exists telegram_id;

-- +goose StatementEnd
