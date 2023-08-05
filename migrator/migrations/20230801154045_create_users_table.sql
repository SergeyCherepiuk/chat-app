-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table users (
    id bigserial primary key,
    first_name varchar(20) not null,
    last_name varchar(20) not null,
    username varchar(30) unique not null,
    password varchar(256) not null,
    description varchar(100) not null,
    profile_picture bytea,
    created_at timestamp not null default current_timestamp
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop table if exists users;