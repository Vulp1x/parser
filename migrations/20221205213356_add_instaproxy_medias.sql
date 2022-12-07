-- +goose Up
-- +goose StatementBegin

CREATE TABLE medias
(
    pk                bigint                                 not null,
    id                text                                   not null,
    dataset_id        uuid references datasets               not null,
    media_type        integer                                not null,
    code              text                                   not null,
    has_more_comments bool                                   not null,
    caption           text                                   not null,
    width             integer                                not null,
    height            integer                                not null,
    like_count        integer                                not null,
    taken_at          integer                                not null,

    created_at        timestamp with time zone DEFAULT now() NOT NULL,
    updated_at        timestamp with time zone DEFAULT now() NOT NULL,
    primary key (pk, dataset_id)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE medias;
-- +goose StatementEnd
