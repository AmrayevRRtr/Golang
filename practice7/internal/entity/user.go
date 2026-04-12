package entity

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"ID";type:uuid;default:uuid_generate_v4()"`
	Username string
	Email    string
	Password string
	Role     string
	Verified bool
}
