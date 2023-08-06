-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table group_chat_messages (
    user_id bigint references users(id) on delete set null,
    chat_id bigint not null references group_chats(id) on delete cascade
) inherits (abstract_message);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop table if exists group_chat_messages;