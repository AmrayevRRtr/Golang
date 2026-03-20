package controllers

import (
	"Practice5/internal/common"
	"Practice5/internal/database"
	"Practice5/internal/helpers"
	"Practice5/internal/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetPaginatedUsers(c *fiber.Ctx) error {
	users := []models.User{}
	perPage := c.Query("per_page", "10")
	sortOrder := c.Query("sort_order", "desc")

	cursor := c.Query("cursor", "")

	limit, err := strconv.ParseInt(perPage, 10, 64)
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if err != nil {
		return c.Status(500).JSON("Invalid per_page option")
	}

	isFirstPage := cursor == ""
	pointsNext := false
	query := database.DB

	if cursor != "" {
		decodedCursor, err := helpers.DecodeCursor(cursor)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}
		pointsNext = decodedCursor["points_next"] == true
		operator, order := GetPaginationOperator(pointsNext, sortOrder)
		whereStr := fmt.Sprintf("(created_at %s ? OR (created_at = ? AND id %s ?))", operator, operator)
		query = query.Where(whereStr, decodedCursor["created_at"], decodedCursor["created_at"], decodedCursor["id"])
		if order != "" {
			sortOrder = order
		}
	}
	query.Order("created_at " + sortOrder).Limit(int(limit) + 1).Find(&users)
	hasPagination := len(users) > int(limit)
	if hasPagination {
		users = users[:limit]
	}
	if !isFirstPage && !pointsNext {
		users = helpers.Reverse(users)
	}
	pageInfo := calculatePagination(isFirstPage, hasPagination, int(limit), users, pointsNext)
	response := common.ResponseDTO{
		Success:    true,
		Data:       users,
		Pagination: pageInfo,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func calculatePagination(isFirstPage bool, hasPagination bool, limit int, users []models.User, pointsNext bool) helpers.PaginationInfo {
	pagination := helpers.PaginationInfo{}
	nextCur := helpers.Cursor{}
	prevCur := helpers.Cursor{}
	if isFirstPage {
		if hasPagination {
			nextCur := helpers.CreateCursor(users[limit-1].ID, users[limit-1].CreatedAt, true)
			pagination = helpers.GeneratePager(nextCur, nil)
		}
	} else {
		if pointsNext {
			if hasPagination {
				nextCur = helpers.CreateCursor(users[limit-1].ID, users[limit-1].CreatedAt, true)
			}
			prevCur = helpers.CreateCursor(users[0].ID, users[0].CreatedAt, false)
			pagination = helpers.GeneratePager(nextCur, prevCur)
		} else {
			nextCur = helpers.CreateCursor(users[limit-1].ID, users[limit-1].CreatedAt, true)
			if hasPagination {
				prevCur = helpers.CreateCursor(users[0].ID, users[0].CreatedAt, false)
			}
			pagination = helpers.GeneratePager(nextCur, prevCur)
		}
	}
	return pagination
}

func GetPaginationOperator(pointsNext bool, sortOrder string) (string, string) {
	if pointsNext && sortOrder == "asc" {
		return ">", ""
	}
	if pointsNext && sortOrder == "desc" {
		return "<", ""
	}
	if !pointsNext && sortOrder == "asc" {
		return "<", "desc"
	}
	if !pointsNext && sortOrder == "desc" {
		return ">", "asc"
	}
	return "", ""
}
