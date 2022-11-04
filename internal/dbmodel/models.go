// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package dbmodel

import (
	"time"

	"github.com/google/uuid"
)

type Blogger struct {
	ID             uuid.UUID  `json:"id"`
	DatasetID      uuid.UUID  `json:"dataset_id"`
	Username       string     `json:"username"`
	UserID         int64      `json:"user_id"`
	FollowersCount int64      `json:"followers_count"`
	IsInitial      bool       `json:"is_initial"`
	CreatedAt      time.Time  `json:"created_at"`
	ParsedAt       time.Time  `json:"parsed_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type Bot struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	SessionID string     `json:"session_id"`
	WorkProxy *string    `json:"work_proxy"`
	IsBlocked bool       `json:"is_blocked"`
	StartedAt *time.Time `json:"started_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type Dataset struct {
	ID               uuid.UUID  `json:"id"`
	TextTemplate     string     `json:"text_template"`
	IncomingAccounts []string   `json:"incoming_accounts"`
	Status           int16      `json:"status"`
	Title            string     `json:"title"`
	CreatedAt        time.Time  `json:"created_at"`
	StartedAt        *time.Time `json:"started_at"`
	StoppedAt        *time.Time `json:"stopped_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}

type Target struct {
	ID        uuid.UUID  `json:"id"`
	DatasetID uuid.UUID  `json:"dataset_id"`
	Username  string     `json:"username"`
	UserID    int64      `json:"user_id"`
	Status    int16      `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
