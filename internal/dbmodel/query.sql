-- name: SaveBotAccounts :copyfrom
insert into bots (username, session_id, work_proxy, is_blocked, started_at)
values ($1, $2, $3, $4, $5);
-- name: SaveTargetUsers :copyfrom
insert into targets (dataset_id, username, user_id)
values ($1, $2, $3);