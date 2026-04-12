package repo

import (
	"fmt"
	"practice7/internal/entity"
	"practice7/internal/pkg/mysql"
)

type UserRepo struct {
	MySQL *mysql.MySQL
}

func NewUserRepo(MySQL *mysql.MySQL) *UserRepo {
	return &UserRepo{MySQL: MySQL}
}

func (u *UserRepo) RegisterUser(user *entity.User) (*entity.User, error) {
	err := u.MySQL.Conn.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) LoginUser(user *entity.LoginUserDTO) (*entity.User, error) {
	var userFromDB entity.User

	if err := u.MySQL.Conn.Where("username = ?", user.Username).First(&userFromDB).Error; err != nil {
		return nil, fmt.Errorf("Username Not Found: %v", err)
	}
	return &userFromDB, nil
}

func (u *UserRepo) GetUserByID(id string) (*entity.User, error) {
	var user entity.User
	if err := u.MySQL.Conn.First(&user, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (u *UserRepo) UpdateUserRole(id, newRole string) error {
	return u.MySQL.Conn.Model(&entity.User{}).Where("id = ?", id).Update("role", newRole).Error

}
