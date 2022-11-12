-- +goose Up
-- +goose StatementBegin

ALTER TABLE targets
    ADD COLUMN full_name text not null default '',
    DROP COLUMN is_business,
    DROP COLUMN followers_count,
    DROP COLUMN followings_count,
    DROP COLUMN contact_phone_number,
    DROP COLUMN public_phone_number,
    DROP COLUMN public_phone_country_code,
    DROP COLUMN city_name,
    DROP COLUMN public_email;

ALTER table bloggers
    ADD COLUMN status smallint not null default 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


-- +goose StatementEnd
