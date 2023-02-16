// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: batch.go

package dbmodel

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

const saveBloggers = `-- name: SaveBloggers :batchexec
insert into bloggers (dataset_id, username, user_id, is_initial, parsed_at, is_private, is_verified, status)
values ($1, $2, $3, false, now(), $4, $5, 'info_saved')
ON CONFLICT (username, dataset_id) DO UPDATE SET parsed_at = excluded.parsed_at
`

type SaveBloggersBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

type SaveBloggersParams struct {
	DatasetID  uuid.UUID `json:"dataset_id"`
	Username   string    `json:"username"`
	UserID     int64     `json:"user_id"`
	IsPrivate  bool      `json:"is_private"`
	IsVerified bool      `json:"is_verified"`
}

func (q *Queries) SaveBloggers(ctx context.Context, arg []SaveBloggersParams) *SaveBloggersBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.DatasetID,
			a.Username,
			a.UserID,
			a.IsPrivate,
			a.IsVerified,
		}
		batch.Queue(saveBloggers, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &SaveBloggersBatchResults{br, len(arg), false}
}

func (b *SaveBloggersBatchResults) Exec(f func(int, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		if b.closed {
			if f != nil {
				f(t, errors.New("batch already closed"))
			}
			continue
		}
		_, err := b.br.Exec()
		if f != nil {
			f(t, err)
		}
	}
}

func (b *SaveBloggersBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}

const saveMedias = `-- name: SaveMedias :batchone
insert into medias(pk, id, dataset_id, media_type, code, has_more_comments, caption, width, height, like_count,
                   taken_at, created_at, updated_at)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, now(), now())
ON CONFLICT (pk, dataset_id) DO UPDATE SET has_more_comments=excluded.has_more_comments,
                                           caption=excluded.caption,
                                           like_count=excluded.like_count,
                                           updated_at=now()
RETURNING pk, id, dataset_id, media_type, code, has_more_comments, caption, width, height, like_count, taken_at, created_at, updated_at
`

type SaveMediasBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

type SaveMediasParams struct {
	Pk              int64     `json:"pk"`
	ID              string    `json:"id"`
	DatasetID       uuid.UUID `json:"dataset_id"`
	MediaType       int32     `json:"media_type"`
	Code            string    `json:"code"`
	HasMoreComments bool      `json:"has_more_comments"`
	Caption         string    `json:"caption"`
	Width           int32     `json:"width"`
	Height          int32     `json:"height"`
	LikeCount       int32     `json:"like_count"`
	TakenAt         int32     `json:"taken_at"`
}

func (q *Queries) SaveMedias(ctx context.Context, arg []SaveMediasParams) *SaveMediasBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.Pk,
			a.ID,
			a.DatasetID,
			a.MediaType,
			a.Code,
			a.HasMoreComments,
			a.Caption,
			a.Width,
			a.Height,
			a.LikeCount,
			a.TakenAt,
		}
		batch.Queue(saveMedias, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &SaveMediasBatchResults{br, len(arg), false}
}

func (b *SaveMediasBatchResults) QueryRow(f func(int, Media, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		var i Media
		if b.closed {
			if f != nil {
				f(t, i, errors.New("batch already closed"))
			}
			continue
		}
		row := b.br.QueryRow()
		err := row.Scan(
			&i.Pk,
			&i.ID,
			&i.DatasetID,
			&i.MediaType,
			&i.Code,
			&i.HasMoreComments,
			&i.Caption,
			&i.Width,
			&i.Height,
			&i.LikeCount,
			&i.TakenAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if f != nil {
			f(t, i, err)
		}
	}
}

func (b *SaveMediasBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}
