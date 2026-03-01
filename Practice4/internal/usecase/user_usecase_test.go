package usecase

import "Practice4/pkg/modules"

type MockUserRepository struct {
	GetUsersFn            func(limit, offset int) ([]modules.User, error)
	GetUserByIDFn         func(id int) (*modules.User, error)
	CreateUserFn          func(user *modules.User) (int64, error)
	UpdateUserFn          func(user *modules.User) error
	DeleteUserFn          func(id int) error
	CreateUserWithAuditFn func(user *modules.User) (int64, error)
	GetUserByEmailFn      func(email string) (*modules.User, error)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*modules.User, error) {
	return m.GetUserByEmailFn(email)
}

func (m *MockUserRepository) GetUsers(limit, offset int) ([]modules.User, error) {
	return m.GetUsersFn(limit, offset)
}

func (m *MockUserRepository) GetUserByID(id int) (*modules.User, error) {
	return m.GetUserByIDFn(id)
}

func (m *MockUserRepository) CreateUser(user *modules.User) (int64, error) {
	return m.CreateUserFn(user)
}

func (m *MockUserRepository) UpdateUser(user *modules.User) error {
	return m.UpdateUserFn(user)
}

func (m *MockUserRepository) DeleteUser(id int) error {
	return m.DeleteUserFn(id)
}

func (m *MockUserRepository) CreateUserWithAudit(user *modules.User) (int64, error) {
	return m.CreateUserWithAuditFn(user)
}
