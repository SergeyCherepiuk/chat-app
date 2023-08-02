-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table users (
    id bigserial primary key,
    first_name varchar(30) not null,
    last_name varchar(30) not null,
    username varchar(20) not null unique,
    password varchar(256) not null
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop table if exists users;