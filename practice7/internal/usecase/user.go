package usecase

import (
	"fmt"
	"os"
	"practice7/internal/entity"
	"practice7/internal/usecase/repo"
	"practice7/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserUseCase struct {
	repo *repo.UserRepo
}

func NewUserUseCase(r *repo.UserRepo) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (u *UserUseCase) RegisterUser(user *entity.User) (*entity.User, string, error) {
	user, err := u.repo.RegisterUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("register user: %w", err)
	}

	sessionID := uuid.New().String()
	return user, sessionID, nil
}

func (u *UserUseCase) LoginUser(user *entity.LoginUserDTO) (string, string, error) {
	userFromRepo, err := u.repo.LoginUser(user)
	if err != nil {
		return "", "", fmt.Errorf("User From Repo: %w", err)
	}
	if !utils.CheckPassword(userFromRepo.Password, user.Password) {
		return "", "", fmt.Errorf("Check Password: %w", err)
	}
	return utils.GenerateTokens(userFromRepo.ID, userFromRepo.Role)
}

func (u *UserUseCase) GetUserByID(id string) (*entity.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *UserUseCase) PromoteUser(id string) error {
	return u.repo.UpdateUserRole(id, "admin")
}

func (u *UserUseCase) RefreshToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	userIDStr := claims["user_id"].(string)

	user, err := u.repo.GetUserByID(userIDStr)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	accessToken, _, err := utils.GenerateTokens(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
