-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table group_chat_tags (
    id bigserial primary key,
    name varchar(20) not null,
    color integer not null,
    chat_id bigint not null references group_chats(id) on delete cascade,
    created_at timestamp not null default current_timestamp
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop table if exists group_chat_tags;