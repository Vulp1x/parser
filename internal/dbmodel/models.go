// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package dbmodel

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BloggerStatus string

const (
	BloggerStatusNew             BloggerStatus = "new"
	BloggerStatusInfoSaved       BloggerStatus = "info_saved"
	BloggerStatusMediasFound     BloggerStatus = "medias_found"
	BloggerStatusAllMediasParsed BloggerStatus = "all_medias_parsed"
	BloggerStatusDone            BloggerStatus = "done"
	BloggerStatusInvalid         BloggerStatus = "invalid"
)

func (e *BloggerStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = BloggerStatus(s)
	case string:
		*e = BloggerStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for BloggerStatus: %T", src)
	}
	return nil
}

type NullBloggerStatus struct {
	BloggerStatus BloggerStatus
	Valid         bool // Valid is true if BloggerStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullBloggerStatus) Scan(value interface{}) error {
	if value == nil {
		ns.BloggerStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.BloggerStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullBloggerStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.BloggerStatus, nil
}

type DatasetType string

const (
	DatasetTypeFollowers        DatasetType = "followers"
	DatasetTypePhoneNumbers     DatasetType = "phone_numbers"
	DatasetTypeLikesAndComments DatasetType = "likes_and_comments"
)

func (e *DatasetType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = DatasetType(s)
	case string:
		*e = DatasetType(s)
	default:
		return fmt.Errorf("unsupported scan type for DatasetType: %T", src)
	}
	return nil
}

type NullDatasetType struct {
	DatasetType DatasetType
	Valid       bool // Valid is true if DatasetType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullDatasetType) Scan(value interface{}) error {
	if value == nil {
		ns.DatasetType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.DatasetType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullDatasetType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.DatasetType, nil
}

type PgqueueStatus string

const (
	PgqueueStatusNew            PgqueueStatus = "new"
	PgqueueStatusMustRetry      PgqueueStatus = "must_retry"
	PgqueueStatusNoAttemptsLeft PgqueueStatus = "no_attempts_left"
	PgqueueStatusCancelled      PgqueueStatus = "cancelled"
	PgqueueStatusSucceeded      PgqueueStatus = "succeeded"
)

func (e *PgqueueStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PgqueueStatus(s)
	case string:
		*e = PgqueueStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PgqueueStatus: %T", src)
	}
	return nil
}

type NullPgqueueStatus struct {
	PgqueueStatus PgqueueStatus
	Valid         bool // Valid is true if PgqueueStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPgqueueStatus) Scan(value interface{}) error {
	if value == nil {
		ns.PgqueueStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PgqueueStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPgqueueStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.PgqueueStatus, nil
}

type Blogger struct {
	ID                     uuid.UUID     `json:"id"`
	DatasetID              uuid.UUID     `json:"dataset_id"`
	Username               string        `json:"username"`
	UserID                 int64         `json:"user_id"`
	FollowersCount         int64         `json:"followers_count"`
	IsInitial              bool          `json:"is_initial"`
	CreatedAt              time.Time     `json:"created_at"`
	ParsedAt               *time.Time    `json:"parsed_at"`
	UpdatedAt              *time.Time    `json:"updated_at"`
	IsCorrect              bool          `json:"is_correct"`
	IsPrivate              bool          `json:"is_private"`
	IsVerified             bool          `json:"is_verified"`
	IsBusiness             bool          `json:"is_business"`
	FollowingsCount        int32         `json:"followings_count"`
	ContactPhoneNumber     *string       `json:"contact_phone_number"`
	PublicPhoneNumber      *string       `json:"public_phone_number"`
	PublicPhoneCountryCode *string       `json:"public_phone_country_code"`
	CityName               *string       `json:"city_name"`
	PublicEmail            *string       `json:"public_email"`
	Status                 BloggerStatus `json:"status"`
}

type Bot struct {
	ID          uuid.UUID  `json:"id"`
	Username    string     `json:"username"`
	SessionID   string     `json:"session_id"`
	Proxy       Proxy      `json:"proxy"`
	IsBlocked   bool       `json:"is_blocked"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	LockedUntil *time.Time `json:"locked_until"`
}

type Dataset struct {
	ID               uuid.UUID     `json:"id"`
	PhoneCode        *int32        `json:"phone_code"`
	Status           datasetStatus `json:"status"`
	Title            string        `json:"title"`
	ManagerID        uuid.UUID     `json:"manager_id"`
	CreatedAt        time.Time     `json:"created_at"`
	StartedAt        *time.Time    `json:"started_at"`
	StoppedAt        *time.Time    `json:"stopped_at"`
	UpdatedAt        *time.Time    `json:"updated_at"`
	DeletedAt        *time.Time    `json:"deleted_at"`
	PostsPerBlogger  int32         `json:"posts_per_blogger"`
	LikedPerPost     int32         `json:"liked_per_post"`
	CommentedPerPost int32         `json:"commented_per_post"`
	Type             DatasetType   `json:"type"`
}

type Media struct {
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
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Pgqueue struct {
	ID              int64         `json:"id"`
	Kind            int16         `json:"kind"`
	Payload         []byte        `json:"payload"`
	ExternalKey     *string       `json:"external_key"`
	Status          PgqueueStatus `json:"status"`
	Messages        []string      `json:"messages"`
	AttemptsLeft    int16         `json:"attempts_left"`
	AttemptsElapsed int16         `json:"attempts_elapsed"`
	DelayedTill     time.Time     `json:"delayed_till"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type Target struct {
	ID         uuid.UUID  `json:"id"`
	Username   string     `json:"username"`
	UserID     int64      `json:"user_id"`
	Status     int16      `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	IsPrivate  bool       `json:"is_private"`
	IsVerified bool       `json:"is_verified"`
	FullName   string     `json:"full_name"`
	MediaPk    int64      `json:"media_pk"`
	DatasetID  uuid.UUID  `json:"dataset_id"`
}
