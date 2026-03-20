package models

import (
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"time"
)

type User struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(50);not null"`
	Email     string    `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Gender    string    `json:"gender" gorm:"type:varchar(10)"`
	BirthDate time.Time `json:"birth_date"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.NewString()
	return
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
