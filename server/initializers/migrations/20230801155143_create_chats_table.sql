-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table chats (
    id bigserial primary key,
    chatter1 bigint references users(id) on delete set null,
    chatter2 bigint references users(id) on delete set null,
    created_at timestamp not null default current_timestamp
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop table if exists chats;