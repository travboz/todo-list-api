package data

import "github.com/travboz/backend-projects/todo-list-api/internal/validator"

type Filters struct {
	Page     int
	PageSize int
	Filters  []string
}

func (p Filters) Limit() int {
	return p.PageSize
}

func (p Filters) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func ValidateFilters(v *validator.Validator, p Filters) {
	// Check that the page and page_size parameters contain sensible values
	v.Check(p.Page > 0, "page", "must be greater than zero")
	v.Check(p.Page <= 10_000, "page", "must be a maximum of 10 million")
	v.Check(p.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(p.PageSize <= 100, "page_size", "must be a maximum of 100")

}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// The calculateMetadata() function calculates the appropriate pagination metadata
// values given the total number of records, current page, and page size values. Note
// that when the last page value is calculated we are dividing two int values, and
// when dividing integer types in Go the result will also be an integer type, with
// the modulus (or remainder) dropped. So, for example, if there were 12 records in total
// and a page size of 5, the last page value would be (12+5-1)/5 = 3.2, which is then
// truncated to 3 by Go.
func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		// Note that we return an empty Metadata struct if there are no records.
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
}
