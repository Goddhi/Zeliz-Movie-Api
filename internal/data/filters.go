package data

import "github.com/goddhi/zeliz-movie/internal/validator"

type Filters struct {
	Page          int
	PageSize      int
	Sort          string
	SortStatelist []string
}

func ValidateFilters(v *validator.Validator, f Filters) {

	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maixmum of 100")

	// check that the sort parameters matches a value in the SortStatelist
	v.Check(validator.In(f.Sort, f.SortStatelist...), "sort", "invalid sort value")

}
