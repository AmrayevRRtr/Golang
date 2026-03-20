package common

import "Practice5/internal/helpers"

type ResponseDTO struct {
	Success    bool                   `json:"success"`
	Data       any                    `json:"data"`
	Pagination helpers.PaginationInfo `json:"pagination"`
}
