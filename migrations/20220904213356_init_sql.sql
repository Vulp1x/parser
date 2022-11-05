-- +goose Up
-- +goose StatementBegin
CREATE
    EXTENSION IF NOT EXISTS pgcrypto;


CREATE TABLE datasets
(
    id         uuid                     not null primary key default gen_random_uuid(),
    phone_code smallint,
    status     smallint                 not null,
    title      text                     not null,
    user_id    uuid                     not null,
    created_at timestamp with time zone not null,
    started_at timestamp with time zone,
    stopped_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

CREATE TABLE bots
(
    id         uuid primary key     default gen_random_uuid(),
    username   text UNIQUE not null,
    session_id text        not null,
    work_proxy jsonb,
    is_blocked bool        not null,
    started_at timestamp,
    created_at timestamp   not null default now(),
    updated_at timestamp,
    deleted_at timestamp
);

create table bloggers
(
    id              uuid primary key   default gen_random_uuid(),
    dataset_id      uuid      not null references datasets,
    username        text      not null,
    user_id         bigint    not null,
    followers_count bigint    not null default -1,
    is_initial      bool      not null, -- является ли блоггер
    created_at      timestamp not null default now(),
    parsed_at       timestamp not null,
    updated_at      timestamp
);


-- таблица с пользователями, которым будет показана реклама
create table targets
(
    id         uuid primary key   default gen_random_uuid(),
    dataset_id uuid      not null references datasets,
    username   text      not null,
    user_id    bigint    not null,
    status     smallint  not null default 0, -- 0 - не показывали рекламу, 1 - пытались показать рекламу, но не получилось, 2-показали рекламу
    created_at timestamp not null default now(),
    updated_at timestamp
);



create unique index target_users_uniq_idx on targets (dataset_id, user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


-- +goose StatementEnd
