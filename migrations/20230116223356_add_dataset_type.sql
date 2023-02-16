-- +goose Up
-- +goose StatementBegin

Set time zone 'Europe/Moscow';

CREATE TYPE dataset_type AS ENUM (
    'followers',
    'phone_numbers',
    'likes_and_comments'
    );

ALTER TABLE datasets
    ADD COLUMN type            dataset_type DEFAULT 'likes_and_comments' NOT NULL,
    ADD COLUMN followers_count integer                                   not null default 0;


create table full_targets
(
    id                            uuid      default gen_random_uuid() not null
        primary key,
    dataset_id                    uuid                                not null references datasets,
    parsed_at                     timestamp default now()             not null,
    username                      text                                not null,
    inst_pk                       bigint                              not null,
    full_name                     text                                not null,
    is_private                    boolean                             not null,
    is_verified                   boolean                             not null,
    is_business                   boolean                             not null,
    is_potential_business         boolean                             not null,
    has_anonymous_profile_picture boolean                             not null,
    biography                     text                                not null,
    external_url                  text                                not null,
    media_count                   integer                             not null,
    follower_count                integer                             not null,
    following_count               integer                             not null,
    category                      text                                not null,
    city_name                     text                                not null,
    contact_phone_number          text                                not null,
    latitude                      double precision                    not null,
    longitude                     double precision                    not null,
    public_email                  text                                not null,
    public_phone_country_code     text                                not null,
    public_phone_number           text                                not null,
    bio_links                     jsonb                               not null,
    whatsapp_number               text                                not null
);

create unique index full_targets_uniq_user_per_dataset
    on full_targets (inst_pk, dataset_id);

ALTER TABLE bloggers
    drop column is_business,
    drop column followings_count,
    drop column contact_phone_number,
    drop column public_phone_number,
    drop column public_phone_country_code,
    drop column city_name,
    drop column public_email;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE dataset_type;
ALTER TABLE datasets
    DROP COLUMN type;
-- +goose StatementEnd
