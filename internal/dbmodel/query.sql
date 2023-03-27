-- name: SelectNow :one
select now()::timestamp with time zone;

-- name: SaveBots :execrows
insert into bots (username, session_id, proxy, is_blocked)
    (select unnest(sqlc.arg(usernames)::text[]),
            unnest(sqlc.arg(session_ids)::text[]),
            unnest(sqlc.arg(proxies)::jsonb[]),
            false)
ON CONFLICT (session_id) DO NOTHING;

-- name: SaveTargetUsers :execrows
insert into targets (username, user_id, full_name, is_private, is_verified, media_pk, dataset_id)
    (select unnest(sqlc.arg(usernames)::text[]),
            unnest(sqlc.arg(user_ids)::bigint[]),
            unnest(sqlc.arg(full_names)::text[]),
            unnest(sqlc.arg(is_private)::bool[]),
            unnest(sqlc.arg(is_verified)::bool[]),
            @media_pk,
            @dataset_id)
ON CONFLICT (user_id, dataset_id) DO UPDATE set updated_at  = now(),
                                                username    = excluded.username,
                                                is_private  = excluded.is_private,
                                                is_verified = excluded.is_verified,
                                                full_name   = excluded.full_name;

-- name: CreateDraftDataset :one
insert into datasets (title, manager_id, type, status, created_at)
VALUES (@title, @manager_id, @type, 1, now())
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
    followers_count    = @followers_count,
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

-- name: FindBloggersForParsing :many
select *
from bloggers
where dataset_id = @dataset_id;
--   AND status = 2
--   AND user_id > 0;

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
where manager_id = @manager_id
order by created_at desc;

-- name: UpdateDatasetStatus :exec
update datasets
set status     = @status,
    updated_at = now()
where id = @id;

-- name: SaveBloggers :batchexec
insert into bloggers (dataset_id, username, user_id, is_initial, parsed_at, is_private, is_verified, status)
values ($1, $2, $3, false, now(), $4, $5, 'info_saved')
ON CONFLICT (username, dataset_id) DO UPDATE SET parsed_at = excluded.parsed_at;

-- name: SetBloggerIsParsed :exec
update bloggers
set is_correct = @is_correct,
    parsed_at  = now()
where id = @id;

-- name: UpdateBlogger :exec
update bloggers
set user_id     = @user_id,
    parsed_at   = @parsed_at,
    is_correct  = @is_correct,
    is_private  = @is_private,
    is_verified = @is_verified,
    status      = 'info_saved'
where id = @id;

-- name: MarkBloggerAsParsed :exec
update bloggers
set status = 'done'
where username = @username
  and dataset_id = @dataset_id;

-- name: MarkBloggerAsSimilarAccountsFound :exec
update bloggers
set status = 'info_saved'
where id = @id;

-- name: GetParsingProgress :one
select (select count(*)
        from bloggers
        where bloggers.dataset_id = @dataset_id
          and status = 'medias_found')                                         as parsed_bloggers_count,
       (select count(*) from bloggers where bloggers.dataset_id = @dataset_id) as total_bloggers,
       (select count(*) from targets where dataset_id = @dataset_id)           as targets_saved_count;

-- name: CountParsedTargets :one
select count(*)
from targets
where dataset_id = @dataset_id
  AND media_pk = @media_pk;

-- name: FindTargetsForDataset :many
select *
from targets
where dataset_id = @dataset_id;

-- name: SaveMedias :batchone
insert into medias(pk, id, dataset_id, media_type, code, has_more_comments, caption, width, height, like_count,
                   taken_at, created_at, updated_at)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, now(), now())
ON CONFLICT (pk, dataset_id) DO UPDATE SET has_more_comments=excluded.has_more_comments,
                                           caption=excluded.caption,
                                           like_count=excluded.like_count,
                                           updated_at=now()
RETURNING *;

-- name: SaveFakeMedia :exec
insert into medias(pk, id, dataset_id, media_type, code, has_more_comments, caption, width, height, like_count,
                   taken_at, created_at, updated_at)
values ($1, $2, $3, -1, '', false, $4, 0, 0, -1, -1, now(), now())
ON CONFLICT (pk, dataset_id) DO UPDATE SET has_more_comments=excluded.has_more_comments,
                                           caption=excluded.caption,
                                           like_count=excluded.like_count,
                                           updated_at=now();

-- name: FindNotReadyBloggers :many
select *
from bloggers
where status = 'new'
  and dataset_id = @dataset_id;

-- name: FindNotParsedBloggers :many
select *
from bloggers
where status = 'new'
  and dataset_id = @dataset_id;

-- name: FindMediaByID :one
select *
from medias
where id = @id;

-- name: SaveFullTarget :exec
insert into full_targets(dataset_id, username, inst_pk, full_name, is_private, is_verified, is_business,
                         is_potential_business, has_anonymous_profile_picture, biography, external_url, media_count,
                         follower_count, following_count, category, city_name, contact_phone_number, latitude,
                         longitude, public_email, public_phone_country_code, public_phone_number, bio_links,
                         whatsapp_number)
VALUES (@dataset_id, @username, @inst_pk, @full_name, @is_private, @is_verified, @is_business,
        @is_potential_business, @has_anonymous_profile_picture, @biography, @external_url, @media_count,
        @follower_count, @following_count, @category, @city_name, @contact_phone_number, @latitude,
        @longitude, @public_email, @public_phone_country_code, @public_phone_number, @bio_links,
        @whatsapp_number)
ON CONFLICT (inst_pk, dataset_id) DO NOTHING;

-- name: FindFullTargetsWithCode :many
select *
from full_targets
where dataset_id = @dataset_id
  and public_phone_country_code = @public_phone_country_code;

-- name: FindFullTargets :many
select *
from full_targets
where dataset_id = @dataset_id;

-- name: SetBloggerStatusToInvalid :exec
update bloggers
set status = 'invalid'
where username = @username
  and dataset_id = @dataset_id;