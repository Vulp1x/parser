// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: query.sql

package dbmodel

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"time"
)

const blockBot = `-- name: BlockBot :exec
update bots
set is_blocked   = true,
    locked_until = null
where id = $1
`

func (q *Queries) BlockBot(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, blockBot, id)
	return err
}

const countAvailableBots = `-- name: CountAvailableBots :one
select count(*)
from bots
where is_blocked = false
`

func (q *Queries) CountAvailableBots(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countAvailableBots)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createDraftDataset = `-- name: CreateDraftDataset :one
insert into datasets (title, manager_id, status, created_at)
VALUES ($1, $2, 1, now())
RETURNING id
`

type CreateDraftDatasetParams struct {
	Title     string    `json:"title"`
	ManagerID uuid.UUID `json:"manager_id"`
}

func (q *Queries) CreateDraftDataset(ctx context.Context, arg CreateDraftDatasetParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createDraftDataset, arg.Title, arg.ManagerID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteBloggersPerDataset = `-- name: DeleteBloggersPerDataset :execresult
delete
from bloggers
where dataset_id = $1
  and is_initial = true
`

func (q *Queries) DeleteBloggersPerDataset(ctx context.Context, datasetID uuid.UUID) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, deleteBloggersPerDataset, datasetID)
}

const findBloggersForDataset = `-- name: FindBloggersForDataset :many
select id, dataset_id, username, user_id, followers_count, is_initial, created_at, parsed_at, updated_at, parsed, is_correct, is_private, is_verified, is_business, followings_count, contact_phone_number, public_phone_number, public_phone_country_code, city_name, public_email, status
from bloggers
where dataset_id = $1
`

func (q *Queries) FindBloggersForDataset(ctx context.Context, datasetID uuid.UUID) ([]Blogger, error) {
	rows, err := q.db.Query(ctx, findBloggersForDataset, datasetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Blogger
	for rows.Next() {
		var i Blogger
		if err := rows.Scan(
			&i.ID,
			&i.DatasetID,
			&i.Username,
			&i.UserID,
			&i.FollowersCount,
			&i.IsInitial,
			&i.CreatedAt,
			&i.ParsedAt,
			&i.UpdatedAt,
			&i.Parsed,
			&i.IsCorrect,
			&i.IsPrivate,
			&i.IsVerified,
			&i.IsBusiness,
			&i.FollowingsCount,
			&i.ContactPhoneNumber,
			&i.PublicPhoneNumber,
			&i.PublicPhoneCountryCode,
			&i.CityName,
			&i.PublicEmail,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findBloggersForParsing = `-- name: FindBloggersForParsing :many
select id, dataset_id, username, user_id, followers_count, is_initial, created_at, parsed_at, updated_at, parsed, is_correct, is_private, is_verified, is_business, followings_count, contact_phone_number, public_phone_number, public_phone_country_code, city_name, public_email, status
from bloggers
where dataset_id = $1
  AND status = 2
  AND user_id > 0
`

func (q *Queries) FindBloggersForParsing(ctx context.Context, datasetID uuid.UUID) ([]Blogger, error) {
	rows, err := q.db.Query(ctx, findBloggersForParsing, datasetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Blogger
	for rows.Next() {
		var i Blogger
		if err := rows.Scan(
			&i.ID,
			&i.DatasetID,
			&i.Username,
			&i.UserID,
			&i.FollowersCount,
			&i.IsInitial,
			&i.CreatedAt,
			&i.ParsedAt,
			&i.UpdatedAt,
			&i.Parsed,
			&i.IsCorrect,
			&i.IsPrivate,
			&i.IsVerified,
			&i.IsBusiness,
			&i.FollowingsCount,
			&i.ContactPhoneNumber,
			&i.PublicPhoneNumber,
			&i.PublicPhoneCountryCode,
			&i.CityName,
			&i.PublicEmail,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findInitialBloggersForDataset = `-- name: FindInitialBloggersForDataset :many
select id, dataset_id, username, user_id, followers_count, is_initial, created_at, parsed_at, updated_at, parsed, is_correct, is_private, is_verified, is_business, followings_count, contact_phone_number, public_phone_number, public_phone_country_code, city_name, public_email, status
from bloggers
where dataset_id = $1
  AND is_initial = true
`

func (q *Queries) FindInitialBloggersForDataset(ctx context.Context, datasetID uuid.UUID) ([]Blogger, error) {
	rows, err := q.db.Query(ctx, findInitialBloggersForDataset, datasetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Blogger
	for rows.Next() {
		var i Blogger
		if err := rows.Scan(
			&i.ID,
			&i.DatasetID,
			&i.Username,
			&i.UserID,
			&i.FollowersCount,
			&i.IsInitial,
			&i.CreatedAt,
			&i.ParsedAt,
			&i.UpdatedAt,
			&i.Parsed,
			&i.IsCorrect,
			&i.IsPrivate,
			&i.IsVerified,
			&i.IsBusiness,
			&i.FollowingsCount,
			&i.ContactPhoneNumber,
			&i.PublicPhoneNumber,
			&i.PublicPhoneCountryCode,
			&i.CityName,
			&i.PublicEmail,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findTargetsForDataset = `-- name: FindTargetsForDataset :many
select id, dataset_id, username, user_id, status, created_at, updated_at, is_private, is_verified, full_name
from targets
where dataset_id = $1
`

func (q *Queries) FindTargetsForDataset(ctx context.Context, datasetID uuid.UUID) ([]Target, error) {
	rows, err := q.db.Query(ctx, findTargetsForDataset, datasetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Target
	for rows.Next() {
		var i Target
		if err := rows.Scan(
			&i.ID,
			&i.DatasetID,
			&i.Username,
			&i.UserID,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.IsPrivate,
			&i.IsVerified,
			&i.FullName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findUserDatasets = `-- name: FindUserDatasets :many
select id, phone_code, status, title, manager_id, created_at, started_at, stopped_at, updated_at, deleted_at, posts_per_blogger, liked_per_post, commented_per_post
from datasets
where manager_id = $1
`

func (q *Queries) FindUserDatasets(ctx context.Context, managerID uuid.UUID) ([]Dataset, error) {
	rows, err := q.db.Query(ctx, findUserDatasets, managerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Dataset
	for rows.Next() {
		var i Dataset
		if err := rows.Scan(
			&i.ID,
			&i.PhoneCode,
			&i.Status,
			&i.Title,
			&i.ManagerID,
			&i.CreatedAt,
			&i.StartedAt,
			&i.StoppedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.PostsPerBlogger,
			&i.LikedPerPost,
			&i.CommentedPerPost,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDatasetByID = `-- name: GetDatasetByID :one
select id, phone_code, status, title, manager_id, created_at, started_at, stopped_at, updated_at, deleted_at, posts_per_blogger, liked_per_post, commented_per_post
from datasets
where id = $1
`

func (q *Queries) GetDatasetByID(ctx context.Context, id uuid.UUID) (Dataset, error) {
	row := q.db.QueryRow(ctx, getDatasetByID, id)
	var i Dataset
	err := row.Scan(
		&i.ID,
		&i.PhoneCode,
		&i.Status,
		&i.Title,
		&i.ManagerID,
		&i.CreatedAt,
		&i.StartedAt,
		&i.StoppedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.PostsPerBlogger,
		&i.LikedPerPost,
		&i.CommentedPerPost,
	)
	return i, err
}

const getParsingProgress = `-- name: GetParsingProgress :one
select (select count(*) from bloggers where bloggers.dataset_id = $1 and status = 3) as parsed_bloggers_count,
       (select count(*) from bloggers where bloggers.dataset_id = $1)                as total_bloggers,
       (select count(*) from targets where targets.dataset_id = $1)                  as targets_saved_coun
`

type GetParsingProgressRow struct {
	ParsedBloggersCount int64 `json:"parsed_bloggers_count"`
	TotalBloggers       int64 `json:"total_bloggers"`
	TargetsSavedCoun    int64 `json:"targets_saved_coun"`
}

func (q *Queries) GetParsingProgress(ctx context.Context, datasetID uuid.UUID) (GetParsingProgressRow, error) {
	row := q.db.QueryRow(ctx, getParsingProgress, datasetID)
	var i GetParsingProgressRow
	err := row.Scan(&i.ParsedBloggersCount, &i.TotalBloggers, &i.TargetsSavedCoun)
	return i, err
}

type InsertInitialBloggersParams struct {
	DatasetID uuid.UUID `json:"dataset_id"`
	Username  string    `json:"username"`
	UserID    int64     `json:"user_id"`
	IsInitial bool      `json:"is_initial"`
}

const lockAvailableBots = `-- name: LockAvailableBots :many
update bots
set locked_until = now() + interval '15m'
where id in (select id
             from bots
             where is_blocked = false
               and (bots.locked_until is null or locked_until < now())
             limit $1)
RETURNING id, username, session_id, proxy, is_blocked, created_at, updated_at, deleted_at, locked_until
`

func (q *Queries) LockAvailableBots(ctx context.Context, limit int32) ([]Bot, error) {
	rows, err := q.db.Query(ctx, lockAvailableBots, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Bot
	for rows.Next() {
		var i Bot
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.SessionID,
			&i.Proxy,
			&i.IsBlocked,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.LockedUntil,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markBloggerAsParsed = `-- name: MarkBloggerAsParsed :exec
update bloggers
set status = 3 -- TargetsParsedBloggerStatus
where id = $1
`

func (q *Queries) MarkBloggerAsParsed(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, markBloggerAsParsed, id)
	return err
}

const markBloggerAsSimilarAccountsFound = `-- name: MarkBloggerAsSimilarAccountsFound :exec
update bloggers
set status = 2 -- TargetsParsedBloggerStatus
where id = $1
`

func (q *Queries) MarkBloggerAsSimilarAccountsFound(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, markBloggerAsSimilarAccountsFound, id)
	return err
}

type SaveBloggersParams struct {
	DatasetID              uuid.UUID     `json:"dataset_id"`
	Username               string        `json:"username"`
	UserID                 int64         `json:"user_id"`
	FollowersCount         int64         `json:"followers_count"`
	IsInitial              bool          `json:"is_initial"`
	ParsedAt               *time.Time    `json:"parsed_at"`
	Parsed                 bool          `json:"parsed"`
	IsPrivate              bool          `json:"is_private"`
	IsVerified             bool          `json:"is_verified"`
	IsBusiness             bool          `json:"is_business"`
	FollowingsCount        int32         `json:"followings_count"`
	ContactPhoneNumber     *string       `json:"contact_phone_number"`
	PublicPhoneNumber      *string       `json:"public_phone_number"`
	PublicPhoneCountryCode *string       `json:"public_phone_country_code"`
	CityName               *string       `json:"city_name"`
	PublicEmail            *string       `json:"public_email"`
	Status                 bloggerStatus `json:"status"`
}

const saveBots = `-- name: SaveBots :execrows
insert into bots (username, session_id, proxy, is_blocked)
    (select unnest($1::text[]),
            unnest($2::text[]),
            unnest($3::jsonb[]),
            false)
ON CONFLICT (session_id) DO NOTHING
`

type SaveBotsParams struct {
	Usernames  []string `json:"usernames"`
	SessionIds []string `json:"session_ids"`
	Proxies    []string `json:"proxies"`
}

func (q *Queries) SaveBots(ctx context.Context, arg SaveBotsParams) (int64, error) {
	result, err := q.db.Exec(ctx, saveBots, arg.Usernames, arg.SessionIds, arg.Proxies)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const saveTargetUsers = `-- name: SaveTargetUsers :execrows
insert into targets (username, user_id, full_name, is_private, is_verified, dataset_id)
    (select unnest($1::text[]),
            unnest($2::bigint[]),
            unnest($3::text[]),
            unnest($4::bool[]),
            unnest($5::bool[]),
            $6)
ON CONFLICT (user_id, dataset_id) DO UPDATE set updated_at  = now(),
                                                username    = excluded.username,
                                                is_private  = excluded.is_private,
                                                is_verified = excluded.is_verified,
                                                full_name   = excluded.full_name
`

type SaveTargetUsersParams struct {
	Usernames  []string  `json:"usernames"`
	UserIds    []int64   `json:"user_ids"`
	FullNames  []string  `json:"full_names"`
	IsPrivate  []bool    `json:"is_private"`
	IsVerified []bool    `json:"is_verified"`
	DatasetID  uuid.UUID `json:"dataset_id"`
}

func (q *Queries) SaveTargetUsers(ctx context.Context, arg SaveTargetUsersParams) (int64, error) {
	result, err := q.db.Exec(ctx, saveTargetUsers,
		arg.Usernames,
		arg.UserIds,
		arg.FullNames,
		arg.IsPrivate,
		arg.IsVerified,
		arg.DatasetID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const setBloggerIsParsed = `-- name: SetBloggerIsParsed :exec
update bloggers
set parsed     = true,
    is_correct = $1,
    parsed_at  = now()
where id = $2
`

type SetBloggerIsParsedParams struct {
	IsCorrect bool      `json:"is_correct"`
	ID        uuid.UUID `json:"id"`
}

func (q *Queries) SetBloggerIsParsed(ctx context.Context, arg SetBloggerIsParsedParams) error {
	_, err := q.db.Exec(ctx, setBloggerIsParsed, arg.IsCorrect, arg.ID)
	return err
}

const unlockBot = `-- name: UnlockBot :exec
update bots
set locked_until = now() + interval '10s'
where id = $1
`

// Чтобы другие запросы смогли опять его использовать
func (q *Queries) UnlockBot(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, unlockBot, id)
	return err
}

const updateBlogger = `-- name: UpdateBlogger :exec
update bloggers
set user_id                   = $1,
    followers_count           = $2,
    parsed_at                 = $3,
    parsed                    = $4,
    is_correct                = $5,
    is_private                = $6,
    is_verified               = $7,
    is_business               = $8,
    followings_count          = $9,
    contact_phone_number      = $10,
    public_phone_number       = $11,
    public_phone_country_code = $12,
    city_name                 = $13,
    public_email              = $14
where id = $15
`

type UpdateBloggerParams struct {
	UserID                 int64      `json:"user_id"`
	FollowersCount         int64      `json:"followers_count"`
	ParsedAt               *time.Time `json:"parsed_at"`
	Parsed                 bool       `json:"parsed"`
	IsCorrect              bool       `json:"is_correct"`
	IsPrivate              bool       `json:"is_private"`
	IsVerified             bool       `json:"is_verified"`
	IsBusiness             bool       `json:"is_business"`
	FollowingsCount        int32      `json:"followings_count"`
	ContactPhoneNumber     *string    `json:"contact_phone_number"`
	PublicPhoneNumber      *string    `json:"public_phone_number"`
	PublicPhoneCountryCode *string    `json:"public_phone_country_code"`
	CityName               *string    `json:"city_name"`
	PublicEmail            *string    `json:"public_email"`
	ID                     uuid.UUID  `json:"id"`
}

func (q *Queries) UpdateBlogger(ctx context.Context, arg UpdateBloggerParams) error {
	_, err := q.db.Exec(ctx, updateBlogger,
		arg.UserID,
		arg.FollowersCount,
		arg.ParsedAt,
		arg.Parsed,
		arg.IsCorrect,
		arg.IsPrivate,
		arg.IsVerified,
		arg.IsBusiness,
		arg.FollowingsCount,
		arg.ContactPhoneNumber,
		arg.PublicPhoneNumber,
		arg.PublicPhoneCountryCode,
		arg.CityName,
		arg.PublicEmail,
		arg.ID,
	)
	return err
}

const updateDataset = `-- name: UpdateDataset :one
update datasets
set phone_code         = $1,
    title              = $2,
    posts_per_blogger  = $3,
    liked_per_post     = $4,
    commented_per_post = $5,
    updated_at         = now()
where id = $6
returning id, phone_code, status, title, manager_id, created_at, started_at, stopped_at, updated_at, deleted_at, posts_per_blogger, liked_per_post, commented_per_post
`

type UpdateDatasetParams struct {
	PhoneCode        *int32    `json:"phone_code"`
	Title            string    `json:"title"`
	PostsPerBlogger  int32     `json:"posts_per_blogger"`
	LikedPerPost     int32     `json:"liked_per_post"`
	CommentedPerPost int32     `json:"commented_per_post"`
	ID               uuid.UUID `json:"id"`
}

func (q *Queries) UpdateDataset(ctx context.Context, arg UpdateDatasetParams) (Dataset, error) {
	row := q.db.QueryRow(ctx, updateDataset,
		arg.PhoneCode,
		arg.Title,
		arg.PostsPerBlogger,
		arg.LikedPerPost,
		arg.CommentedPerPost,
		arg.ID,
	)
	var i Dataset
	err := row.Scan(
		&i.ID,
		&i.PhoneCode,
		&i.Status,
		&i.Title,
		&i.ManagerID,
		&i.CreatedAt,
		&i.StartedAt,
		&i.StoppedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.PostsPerBlogger,
		&i.LikedPerPost,
		&i.CommentedPerPost,
	)
	return i, err
}

const updateDatasetStatus = `-- name: UpdateDatasetStatus :exec
update datasets
set status     = $1,
    updated_at = now()
where id = $2
`

type UpdateDatasetStatusParams struct {
	Status datasetStatus `json:"status"`
	ID     uuid.UUID     `json:"id"`
}

func (q *Queries) UpdateDatasetStatus(ctx context.Context, arg UpdateDatasetStatusParams) error {
	_, err := q.db.Exec(ctx, updateDatasetStatus, arg.Status, arg.ID)
	return err
}
