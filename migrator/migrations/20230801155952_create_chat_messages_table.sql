-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table chat_messages (
    message_from bigint not null references users(id) on delete cascade,
    message_to bigint not null references users(id) on delete cascade
) inherits (abstract_message);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop table if exists chat_messages;