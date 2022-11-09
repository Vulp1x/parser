-- +goose Up
-- +goose StatementBegin
ALTER TABLE bots
    DROP COLUMN started_at,
    ADD COLUMN locked_until timestamp;

CREATE UNIQUE INDEX on bots (session_id);

ALTER TABLE bloggers
    ADD COLUMN parsed                    boolean not null default false,
    ADD COLUMN is_private                boolean not null default false,
    ADD COLUMN is_verified               boolean not null default false,
    ADD COLUMN is_business               boolean not null default false,
    ADD COLUMN followings_count          integer not null default -1,
    ADD COLUMN contact_phone_number      text,
    ADD COLUMN public_phone_number       text,
    ADD COLUMN public_phone_country_code text,
    ADD COLUMN city_name                 text,
    ADD COLUMN public_email              text;


ALTER TABLE targets
    ADD COLUMN is_private                boolean not null default false,
    ADD COLUMN is_verified               boolean not null default false,
    ADD COLUMN is_business               boolean not null default false,
    ADD COLUMN followers_count           integer not null default -1,
    ADD COLUMN followings_count          integer not null default -1,
    ADD COLUMN contact_phone_number      text,
    ADD COLUMN public_phone_number       text,
    ADD COLUMN public_phone_country_code text,
    ADD COLUMN city_name                 text,
    ADD COLUMN public_email              text;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


-- +goose StatementEnd
