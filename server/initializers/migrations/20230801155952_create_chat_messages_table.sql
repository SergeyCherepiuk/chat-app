-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table chat_messages (
    id bigserial primary key,
    user_id bigint references users(id) on delete set null,
    chat_id bigint not null references chats(id) on delete cascade,
    message text not null,
    is_edited boolean not null,
    created_at timestamp not null default current_timestamp
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop table if exists chat_messages;