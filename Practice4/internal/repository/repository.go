package repository

import "Practice4/pkg/modules"

type UserRepository interface {
	GetUsers(limit, offset int) ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	CreateUser(user *modules.User) (int64, error)
	UpdateUser(user *modules.User) error
	DeleteUser(id int) error
	CreateUserWithAudit(user *modules.User) (int64, error)
	GetUserByEmail(email string) (*modules.User, error)
}

type Repositories struct {
	UserRepository
}

func NewRepositories(userRepo UserRepository) *Repositories {
	return &Repositories{
		UserRepository: userRepo,
	}
}
