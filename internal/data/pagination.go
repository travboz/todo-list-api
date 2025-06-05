package data

import "github.com/travboz/backend-projects/todo-list-api/internal/validator"

type Pagination struct {
	Page     int
	PageSize int
	Filters  []string
}

func (p Pagination) Limit() int {
	return p.PageSize
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func ValidatePagination(v *validator.Validator, p Pagination) {
	// Check that the page and page_size parameters contain sensible values
	v.Check(p.Page > 0, "page", "must be greater than zero")
	v.Check(p.Page <= 10_000, "page", "must be a maximum of 10 million")
	v.Check(p.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(p.PageSize <= 100, "page_size", "must be a maximum of 100")

}
