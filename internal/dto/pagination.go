package dto

import "math"

type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	start       int `json:"-"` // index
	end         int `json:"-"` // index, exclude
}

func (p Pagination) Start() int {
	return p.start
}

func (p Pagination) End() int {
	return p.end
}

func NewPagination(page int, limit int, total int) Pagination {
	if total <= 0 {
		return Pagination{
			CurrentPage: 1,
			TotalPage:   0,
		}
	}

	currentPage := page

	if page <= 0 {
		currentPage = 1
	}

	if limit <= 0 {
		limit = total
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	if currentPage > totalPages {
		currentPage = totalPages
	}

	start := (currentPage - 1) * limit
	end := start + limit
	if end > total {
		end = total
	}

	return Pagination{
		CurrentPage: currentPage,
		TotalPage:   totalPages,
		start:       start,
		end:         end,
	}

}
