-- name: SaveBotAccounts :copyfrom
insert into bots (username, session_id, proxy, is_blocked, started_at)
values ($1, $2, $3, $4, $5);
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
where dataset_id = $1
  and is_initial = true;

-- name: InsertInitialBloggers :copyfrom
insert into bloggers(dataset_id, username, user_id, is_initial)
VALUES ($1, $2, $3, $4);

-- name: FindBloggersForDataset :many
select *
from bloggers
where dataset_id = $1;

-- name: FindUserDatasets :many
select *
from datasets
where manager_id = $1;