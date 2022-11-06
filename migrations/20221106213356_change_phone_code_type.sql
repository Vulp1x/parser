-- +goose Up
-- +goose StatementBegin
ALTER TABLE datasets
    ALTER COLUMN phone_code SET DATA TYPE integer,
    ADD COLUMN posts_per_blogger  integer not null default 0,
    ADD COLUMN liked_per_post     integer not null default 0,
    ADD COLUMN commented_per_post integer not null default 0;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


-- +goose StatementEnd
