-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table abstract_message (
    id bigserial primary key,
    message text not null,
    is_edited boolean not null,
    created_at timestamp not null default current_timestamp
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop table if exists abstract_message;