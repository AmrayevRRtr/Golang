package usecase

import (
	"Practice4/internal/repository"
	"Practice4/pkg/modules"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func (u *UserUsecase) GetUserByEmail(email string) (*modules.User, error) {
	return u.repo.GetUserByEmail(email)
}

func (u *UserUsecase) CreateUserWithAudit(user *modules.User) (int64, error) {
	return u.repo.CreateUserWithAudit(user)
}

func NewUserUsecase(repo repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) GetUsers(limit, offset int) ([]modules.User, error) {
	return u.repo.GetUsers(limit, offset)
}

func (u *UserUsecase) GetUserByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *UserUsecase) CreateUser(user *modules.User) (int64, error) {

	hashed, err := hashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	user.Password = hashed
	return u.repo.CreateUser(user)
}

func (u *UserUsecase) UpdateUser(user *modules.User) error {
	return u.repo.UpdateUser(user)
}

func (u *UserUsecase) DeleteUser(id int) error {
	return u.repo.DeleteUser(id)
}
