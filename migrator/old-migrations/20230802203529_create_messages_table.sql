-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table messages (
    id bigserial primary key,
    message text not null,
    sent_at timestamp not null,
    user_id bigint references users(id),
    chat_id bigint references chats(id) not null
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop table if exists messages;