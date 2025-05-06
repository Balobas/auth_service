-- +goose Up
-- +goose StatementBegin

alter table user_permissions alter column permissions type varchar(100)[] using permissions::text::varchar(100)[];
drop type user_permission;

create table permissions (
    key varchar(100) PRIMARY KEY,
    description varchar(200)
);

insert into permissions (key, description) values 
('not_verified', 'Пользователь не верифицирован. Минимальные возможные права'),
('base', 'Стандартные права пользователя');

create table deleting_permissions (
    key varchar(100) primary key
);

insert into config(key, value) values
('remove_permissions_interval', '"40s"'),
('users_limit_on_remove_permission', '50');

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
delete from config where key in ('remove_permissions_interval', 'users_limit_on_remove_permission');
drop table deleting_permissions;
drop table permissions;
create type user_permission as enum ('base', 'not_verified');
alter table user_permissions alter column permissions type user_permission[] using permissions::text::user_permission[];

-- +goose StatementEnd
