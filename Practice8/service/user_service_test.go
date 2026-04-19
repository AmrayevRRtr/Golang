package service

import (
	"Practice8/repository"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	ctx := context.Background()

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	result, err := userService.GetUserByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}
func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	ctx := context.Background()

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := userService.CreateUser(ctx, user)
	assert.NoError(t, err)
}

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	ctx := context.Background()

	user := &repository.User{ID: 1, Name: "Test"}

	t.Run("user already exists", func(t *testing.T) {
		mockRepo.EXPECT().GetByEmail("test@mail.com").Return(user, nil)

		err := service.RegisterUser(ctx, user, "test@mail.com")

		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("new user success", func(t *testing.T) {
		mockRepo.EXPECT().GetByEmail("new@mail.com").Return(nil, nil)
		mockRepo.EXPECT().CreateUser(user).Return(nil)

		ctx := context.Background()

		err := service.RegisterUser(ctx, user, "new@mail.com")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().GetByEmail("err@mail.com").Return(nil, fmt.Errorf("db error"))

		err := service.RegisterUser(ctx, user, "err@mail.com")

		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}

func TestUpdateUserName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	ctx := context.Background()

	t.Run("empty name", func(t *testing.T) {
		err := service.UpdateUserName(ctx, 1, "")

		if err == nil {
			t.Errorf("expected error")
		}
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(1).Return(nil, fmt.Errorf("not found"))

		err := service.UpdateUserName(ctx, 1, "New")

		if err == nil {
			t.Errorf("expected error")
		}
	})

	t.Run("success", func(t *testing.T) {
		user := &repository.User{ID: 1, Name: "Old"}

		mockRepo.EXPECT().GetUserByID(1).Return(user, nil)
		mockRepo.EXPECT().UpdateUser(user).Return(nil)

		err := service.UpdateUserName(ctx, 1, "New")

		if err != nil {
			t.Errorf("unexpected error")
		}

		if user.Name != "New" {
			t.Errorf("name not updated")
		}
	})
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	ctx := context.Background()

	t.Run("delete admin", func(t *testing.T) {
		err := service.DeleteUser(ctx, 1)

		if err == nil {
			t.Errorf("expected error")
		}
	})

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().DeleteUser(2).Return(nil)

		err := service.DeleteUser(ctx, 2)

		if err != nil {
			t.Errorf("unexpected error")
		}
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.EXPECT().DeleteUser(3).Return(fmt.Errorf("db error"))

		err := service.DeleteUser(ctx, 3)

		if err == nil {
			t.Errorf("expected error")
		}
	})
}
