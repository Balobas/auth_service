-- +goose Up
-- +goose StatementBegin
create table services (
    uid uuid primary key,
    name varchar not null,
    domain varchar not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table services_credentials (
    service_uid uuid references services(uid) on delete cascade,
    h_password varchar
);

create table services_roles (
    service_uid uuid references services(uid) on delete cascade,
    role varchar(50) references roles(role) on delete cascade,
    primary key(service_uid, role)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
