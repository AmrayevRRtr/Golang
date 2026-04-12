package usecase

import "practice7/internal/entity"

type (
	UserInterface interface {
		LoginUser(user *entity.LoginUserDTO) (string, string, error)
		RegisterUser(user *entity.User) (*entity.User, string, error)
		GetUserByID(id string) (*entity.User, error)
		PromoteUser(id string) error
		RefreshToken(tokenStr string) (string, error)
	}
)
