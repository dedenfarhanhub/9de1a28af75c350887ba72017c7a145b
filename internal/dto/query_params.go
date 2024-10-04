package dto

// QueryParams struct
type QueryParams struct {
	Search    string `json:"search" binding:"omitempty"`
	Page      int    `json:"page" binding:"required"`
	PageSize  int    `json:"page_size" binding:"required"`
	SortBy    string `json:"sort_by" binding:"omitempty"`
	SortOrder string `json:"sort_order" binding:"omitempty"`
}
