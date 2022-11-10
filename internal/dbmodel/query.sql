-- name: SaveBots :execrows
insert into bots (username, session_id, proxy, is_blocked)
    (select unnest(sqlc.arg(usernames)::text[]),
            unnest(sqlc.arg(session_ids)::text[]),
            unnest(sqlc.arg(proxies)::jsonb[]),
            false)
ON CONFLICT (session_id) DO NOTHING;

-- name: SaveTargetUsers :copyfrom
insert into targets (dataset_id, username, user_id)
values ($1, $2, $3);

-- name: CreateDraftDataset :one
insert into datasets (title, manager_id, status, created_at)
VALUES (@title, @manager_id, 1, now())
RETURNING id;

-- name: GetDatasetByID :one
select *
from datasets
where id = @id;

-- name: UpdateDataset :one
update datasets
set phone_code         = @phone_code,
    title              = @title,
    posts_per_blogger  = @posts_per_blogger,
    liked_per_post     = @liked_per_post,
    commented_per_post = @commented_per_post,
    updated_at         = now()
where id = @id
returning *;

-- name: DeleteBloggersPerDataset :execresult
delete
from bloggers
where dataset_id = @dataset_id
  and is_initial = true;

-- name: InsertInitialBloggers :copyfrom
insert into bloggers(dataset_id, username, user_id, is_initial)
VALUES (?, ?, ?, ?);

-- name: FindBloggersForDataset :many
select *
from bloggers
where dataset_id = @dataset_id;

-- name: FindInitialBloggersForDataset :many
select *
from bloggers
where dataset_id = @dataset_id
  AND is_initial = true;

-- name: LockAvailableBots :many
update bots
set locked_until = now() + interval '15m'
where id in (select id
             from bots
             where is_blocked = false
               and (bots.locked_until is null or locked_until < now())
             limit $1)
RETURNING *;

-- Чтобы другие запросы смогли опять его использовать
-- name: UnlockBot :exec
update bots
set locked_until = now() + interval '10s'
where id = @id;

-- name: BlockBot :exec
update bots
set is_blocked   = true,
    locked_until = null
where id = @id;

-- name: CountAvailableBots :one
select count(*)
from bots
where is_blocked = false;

-- name: FindUserDatasets :many
select *
from datasets
where manager_id = @manager_id;

-- name: UpdateDatasetStatus :exec
update datasets
set status     = @status,
    updated_at = now()
where id = @id;

-- name: SaveBloggers :copyfrom
insert into bloggers (dataset_id, username, user_id, followers_count, is_initial, parsed_at,
                      parsed, is_private, is_verified, is_business, followings_count, contact_phone_number,
                      public_phone_number, public_phone_country_code, city_name, public_email)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);

-- name: SetBloggerIsParsed :exec
update bloggers
set parsed     = true,
    is_correct = @is_correct,
    parsed_at  = now()
where id = @id;

-- name: UpdateBlogger :exec
update bloggers
set user_id                   = @user_id,
    followers_count           = @followers_count,
    parsed_at                 = @parsed_at,
    parsed                    = @parsed,
    is_correct                = @is_correct,
    is_private                = @is_private,
    is_verified               = @is_verified,
    is_business               = @is_business,
    followings_count          = @followings_count,
    contact_phone_number      = @contact_phone_number,
    public_phone_number       = @public_phone_number,
    public_phone_country_code = @public_phone_country_code,
    city_name                 = @city_name,
    public_email              = @public_email
where id = @id;