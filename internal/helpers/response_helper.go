package helpers

import "github.com/dedenfarhanhub/blog-service/internal/dto"

// NewErrorResponse creates a new BaseResponse for errors.
func NewErrorResponse(code int, message string) *dto.BaseResponse {
	return &dto.BaseResponse{
		Code:    code,
		Status:  "ERROR",
		Message: message,
		Data:    nil,
	}
}

// NewSuccessResponse creates a new BaseResponse for successful responses.
func NewSuccessResponse(data interface{}) *dto.BaseResponse {
	return &dto.BaseResponse{
		Code:    200,
		Status:  "SUCCESS",
		Message: "SUCCESS",
		Data:    data,
	}
}

// NewSuccessResponsePagination creates a new BaseResponse for successful responses.
func NewSuccessResponsePagination(item interface{}, totalCount int64) *dto.BaseResponse {
	return &dto.BaseResponse{
		Code:    200,
		Status:  "SUCCESS",
		Message: "SUCCESS",
		Data: dto.PaginationResponse{
			Items:      item,
			TotalCount: totalCount,
		},
	}
}
