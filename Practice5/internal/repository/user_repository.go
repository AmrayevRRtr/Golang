package repository

import (
	"Practice5/internal/models"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetPaginatedUsers(page int, pageSize int) (models.PaginatedResponse, error) {
	var users []models.User
	offset := (page - 1) * pageSize

	var totalCount int
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`

	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return models.PaginatedResponse{}, err
	}

	query := `SELECT id, name, email, gender, birth_date, deleted_at 
            FROM users 
            WHERE deleted_at IS NULL
            ORDER BY id LIMIT ? OFFSET ?`
	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return models.PaginatedResponse{}, err
	}

	defer rows.Close()
	for rows.Next() {
		var u models.User

		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate, &u.DeletedAt); err != nil {
			return models.PaginatedResponse{}, err
		}
		users = append(users, u)
	}

	return models.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (r *Repository) CreateUser(user *models.User) (int64, error) {
	result, err := r.db.Exec("INSERT INTO users (name, email, gender, birth_date) VALUES (?, ?, ?, ?)",
		user.Name, user.Email, user.Gender, user.BirthDate)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) GetUserByID(id int) (*models.User, error) {
	var user models.User

	err := r.db.Get(&user, "SELECT * FROM users WHERE id=? and deleted_at IS NULL", id)
	if err != nil {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	return &user, nil
}

func (r *Repository) UpdateUser(user *models.User) error {
	result, err := r.db.Exec("UPDATE users SET name=?, email=?, gender=?, birth_date=? WHERE id=?",
		user.Name, user.Email, user.Gender, user.BirthDate, user.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user with id %d not found", user.ID) // ИСПРАВЛЕНО: теперь возвращает ошибку
	}
	return nil
}

func (r *Repository) DeleteUser(id int) error {
	result, err := r.db.Exec("UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id=? AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found or already deleted")
	}
	return nil
}

func (r *Repository) queryRowsWithFilter(ctx context.Context, query string, filter *models.UserFilter) (*sql.Rows, error) {
	var filterValues []any

	query += " WHERE 1=1"

	if !filter.ID.IsZero() {
		query += " AND id = ?"
		filterValues = append(filterValues, filter.ID.Int64)
	}

	if !filter.Name.IsZero() {
		query += " AND name LIKE ?"
		filterValues = append(filterValues, "%"+filter.Name.String+"%")
	}

	if !filter.Email.IsZero() {
		query += " AND email LIKE ?"
		filterValues = append(filterValues, "%"+filter.Email.String+"%")
	}

	if !filter.Gender.IsZero() {
		query += " AND gender = ?"
		filterValues = append(filterValues, filter.Gender.String)
	}

	if !filter.BirthDate.IsZero() {
		query += " AND birth_date = ?"
		filterValues = append(filterValues, filter.BirthDate.Time)
	}

	if filter.LastSeenId > 0 {
		query += " AND id < ?"
		filterValues = append(filterValues, filter.LastSeenId)
	}

	sortField := "id"
	if !filter.SortBy.IsZero() {
		switch filter.SortBy.String {
		case "id", "name", "email", "gender", "birth_date":
			sortField = filter.SortBy.String
		}
	}

	if filter.ShowDeleted {
		query += " AND deleted_at IS NOT NULL"
	} else {
		query += " AND deleted_at IS NULL"
	}

	order := "DESC"
	if !filter.Order.IsZero() && filter.Order.String == "asc" {
		order = "ASC"
	}

	query += " ORDER BY " + sortField + " " + order + " LIMIT ?"
	limit := filter.Limit
	if limit == 0 {
		limit = 10
	}
	filterValues = append(filterValues, limit)

	return r.db.QueryContext(ctx, query, filterValues...)
}

func (r *Repository) ListByFilter(ctx context.Context, filter *models.UserFilter) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	querySelect := `SELECT id, name, email, gender, birth_date, deleted_at FROM users`

	rows, err := r.queryRowsWithFilter(ctx, querySelect, filter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var arr []*models.User
	for rows.Next() {
		u := new(models.User)
		// Теперь тут 6 колонок и 6 аргументов
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Gender,
			&u.BirthDate,
			&u.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		arr = append(arr, u)
	}
	return arr, rows.Err()
}

func (r *Repository) CountByFilter(ctx context.Context, filter *models.UserFilter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var count uint64
	var filterValues []any

	query := `SELECT COUNT(id) FROM users WHERE 1=1`

	if !filter.ID.IsZero() {
		query += " AND id = ?"
		filterValues = append(filterValues, filter.ID.Int64)
	}
	if !filter.Name.IsZero() {
		query += " AND name LIKE ?"
		filterValues = append(filterValues, "%"+filter.Name.String+"%") // Исправлено на LIKE с %
	}
	if !filter.Email.IsZero() {
		query += " AND email LIKE ?"
		filterValues = append(filterValues, "%"+filter.Email.String+"%") // Исправлено на LIKE с %
	}
	if !filter.Gender.IsZero() {
		query += " AND gender = ?"
		filterValues = append(filterValues, filter.Gender.String)
	}
	if !filter.BirthDate.IsZero() {
		query += " AND birth_date = ?"
		filterValues = append(filterValues, filter.BirthDate.Time)
	}

	err := r.db.QueryRowContext(ctx, query, filterValues...).Scan(&count)
	return count, err
}

func (r *Repository) GetCommonFriends(ctx context.Context, user1ID, user2ID int64) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
       SELECT u.id, u.name, u.email, u.gender, u.birth_date
       FROM users u
       JOIN user_friends f1 ON u.id = f1.friend_id
       JOIN user_friends f2 ON u.id = f2.friend_id
       WHERE f1.user_id = ? and f2.user_id = ?`

	rows, err := r.db.QueryContext(ctx, query, user1ID, user2ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []*models.User
	for rows.Next() {
		u := new(models.User)
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Gender,
			&u.BirthDate,
		)
		if err != nil {
			return nil, err
		}
		friends = append(friends, u)
	}
	return friends, rows.Err()
}
