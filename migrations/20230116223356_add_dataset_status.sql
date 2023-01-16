-- +goose Up
-- +goose StatementBegin

Set time zone 'Europe/Moscow';

CREATE TYPE dataset_type AS ENUM (
    'followers',
    'phone_numbers',
    'likes_and_comments'
    );

ALTER TABLE datasets
    ADD COLUMN type dataset_type DEFAULT 'likes_and_comments' NOT NULL;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE dataset_type;
ALTER TABLE datasets
    DROP COLUMN type;
-- +goose StatementEnd
