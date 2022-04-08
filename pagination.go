package db

// PaginationParams should've used to pass pagination parameters to data layer.
type PaginationParams struct {
	Limit  int
	SortBy string
	Page   int
}
