-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table group_chats (
    id bigserial primary key,
    name varchar(40) not null,
    creator_id bigint references users(id) on delete set null,
    created_at timestamp not null default current_timestamp
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop table if exists group_chats;