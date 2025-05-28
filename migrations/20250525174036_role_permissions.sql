-- +goose Up
-- +goose StatementBegin

create table roles (
    role varchar(50) primary key,
    description varchar(200)
);

create table user_roles (
    user_uid uuid references users(uid) on delete cascade,
    role varchar(50) references roles(role) on delete cascade
);

insert into roles (role, description) values ('admin', 'god of the system');

create table roles_permissions (
    role varchar(50) references roles(role) on delete cascade,
    permission varchar(100) references permissions(key) on delete cascade
);

create table resources_permissions (
    uri varchar(100),
    method varchar(10),
    permission varchar(100) references permissions(key),
    primary key (uri, method)
);

alter table users drop column role;
drop type user_role;

drop table user_permissions;

alter table users add column is_verified boolean not null default false;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

create table user_permissions (
    user_uid uuid primary key,
    permissions varchar(100)[],
    foreign key(user_uid) references users(uid) on delete cascade
);

create type user_role as enum('user', 'admin');

alter table users add column role user_role;
drop table resources_permissions;
drop table roles_permissions;
drop table roles;

-- +goose StatementEnd
