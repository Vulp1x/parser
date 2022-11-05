-- name: SaveBotAccounts :copyfrom
insert into bots (username, session_id, work_proxy, is_blocked, started_at)
values ($1, $2, $3, $4, $5);
-- name: SaveTargetUsers :copyfrom
insert into targets (dataset_id, username, user_id)
values ($1, $2, $3);

-- name: CreateDraftDataset :one
insert into datasets (status, title, user_id, status, created_at)
VALUES (@status, @title, @user_id, 1, now())
RETURNING id;

-- name: GetDatasetByID :one
select *
from datasets where id = @id;