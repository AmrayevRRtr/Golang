package repository

import "Practice4/pkg/modules"

type UserRepository interface {
	GetUsers(limit, offset int) ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	CreateUser(user *modules.User) (int64, error)
	UpdateUser(user *modules.User) error
	DeleteUser(id int) error
	CreateUserWithAudit(user *modules.User) (int64, error)
}

type Repositories struct {
	UserRepository
}

func (*Repositories) CreateUserWithAudit(user *modules.User) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func NewRepositories(userRepo UserRepository) *Repositories {
	return &Repositories{
		UserRepository: userRepo,
	}
}
