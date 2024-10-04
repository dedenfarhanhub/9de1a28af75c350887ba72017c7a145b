package dto

// BaseResponse represents the standard structure for API responses.
type BaseResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginationResponse represents the structure for paginated responses.
type PaginationResponse struct {
	Items      interface{} `json:"items,omitempty"`
	TotalCount int64       `json:"total_count"`
}
