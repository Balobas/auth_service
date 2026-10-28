-- +goose Up
-- +goose StatementBegin

-- Таблица дневной активности для DAU/MAU. Одна строка на (пользователь, день).
-- Пишется апсертом на каждом авторизованном запросе (через VerifyAuth), но
-- клиент дедуплицирует в памяти → максимум одна запись на юзера в день.
create table user_activity (
    user_uid    uuid not null,
    active_date date not null,
    primary key (user_uid, active_date)
);

create index idx_user_activity_date on user_activity (active_date);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table user_activity;
-- +goose StatementEnd
