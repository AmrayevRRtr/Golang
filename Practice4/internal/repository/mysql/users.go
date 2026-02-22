package mysql

import (
	"Practice4/pkg/modules"
	"fmt"
)

type UserRepo struct {
	db *Dialect
}

func NewUserRepository(db *Dialect) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetUsers(limit, offset int) ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.Select(&users, "SELECT * FROM users WHERE deleted_at IS NULL LIMIT ? OFFSET ?", limit, offset)

	return users, err
}

func (r *UserRepo) GetUserByID(id int) (*modules.User, error) {
	var user modules.User
	err := r.db.DB.Select(&user, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	return &user, nil
}

func (r *UserRepo) CreateUser(user *modules.User) (int64, error) {
	result, err := r.db.DB.Exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
		user.Name, user.Email, user.Age)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *UserRepo) UpdateUser(user *modules.User) error {
	result, err := r.db.DB.Exec("UPDATE users SET name=?, email=?, age=? WHERE id=?",
		user.Name, user.Email, user.Age, user.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user with id %d not found", user.ID)
	}
	return nil
}

func (r *UserRepo) DeleteUser(id int) error {
	result, err := r.db.DB.Exec("UPDATE users SET deleted_at = NOW() WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}
	return nil
}

func (r *UserRepo) CreateUserWithAudit(user *modules.User) (int64, error) {
	tx, err := r.db.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transacrtion: %v", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	res, err := tx.Exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)", user.Name, user.Email, user.Age)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert user: %v", err)
	}

	userID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to get ladt insert id: %v", err)
	}

	_, err = tx.Exec("INSERT INTO audit_logs (action) VALUES (?)", fmt.Sprintf("Created user with ID %d", userID))
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert audit log: %v", err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}
	return userID, nil
}
