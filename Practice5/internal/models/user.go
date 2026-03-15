package models

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
	DeletedAt null.Time `json:"deleted_at"`
}

type UserFilter struct {
	ID          null.Int    `json:"id"`
	Name        null.String `json:"name"`
	Email       null.String `json:"email"`
	Gender      null.String `json:"gender"`
	BirthDate   null.Time   `json:"birth_date"`
	ShowDeleted bool        `json:"show_deleted"`

	SortBy null.String `json:"sort_by"`
	Order  null.String `json:"order"`

	LastSeenId uint64 `json:"last_seen_id"`
	Limit      uint64 `json:"limit"`
	Offset     uint64 `json:"offset"`
}
