-- +goose Up
-- +goose StatementBegin

Set time zone 'Europe/Moscow';

CREATE TYPE blogger_status AS ENUM (
    'new',
    'info_saved',
    'medias_found',
    'all_medias_parsed',
    'done',
    'invalid'
    );

ALTER TABLE bloggers
    DROP COLUMN status,
    DROP COLUMN parsed,
    ADD COLUMN status blogger_status DEFAULT 'new' NOT NULL;

CREATE UNIQUE INDEX uniq_bloggers_per_dataset ON bloggers (username, dataset_id);

ALTER TABLE targets
    ADD COLUMN media_pk                 bigint not null,
    ADD COLUMN IF NOT EXISTS dataset_id uuid   not null,
    ADD CONSTRAINT medias_fk FOREIGN KEY (media_pk, dataset_id) REFERENCES medias (pk, dataset_id);

CREATE UNIQUE INDEX targets_uniq_user_per_dataset ON targets (user_id, dataset_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE pgqueue_status;
DROP TABLE pgqueue;
-- +goose StatementEnd
