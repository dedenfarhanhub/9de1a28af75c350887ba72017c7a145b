package controllers

import (
	"github.com/dedenfarhanhub/blog-service/internal/dto"
	"github.com/dedenfarhanhub/blog-service/internal/helpers"
	"github.com/dedenfarhanhub/blog-service/internal/services"
	"html"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommentController struct
type CommentController struct {
	commentService services.CommentService
}

// NewCommentController initializes comment controller
func NewCommentController(commentService services.CommentService) *CommentController {
	return &CommentController{commentService: commentService}
}

// Create godoc
// @Summary Create a new comment
// @Description Add a new comment to a post
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param comment body dto.CommentRequest true "Comment details"
// @Success 200 {object} dto.BaseResponse{data=dto.CommentResponse}
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /posts/{id}/comments [post]
func (c *CommentController) Create(ctx *gin.Context) {
	postIDParam := ctx.Param("id")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil || postID <= 0 {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "Invalid Post ID"))
		return
	}

	var commentDto dto.CommentRequest
	if err := ctx.ShouldBindJSON(&commentDto); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// Sanitasi input untuk mencegah XSS
	commentDto.Content = html.EscapeString(commentDto.Content)

	commentResponses, err := c.commentService.Create(uint(postID), &commentDto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewSuccessResponse(commentResponses))
}

// GetAllByPostID godoc
// @Summary Get all comments by post ID
// @Description Retrieve all comments associated with a specific post
// @Tags Comments
// @Produce json
// @Param id path int true "Post ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Number of comments per page"
// @Param search query string false "Search by comment content"
// @Param sort_by query string false "Sort by field"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Success 200 {object} dto.BaseResponse{data=dto.PaginationResponse{items=[]dto.CommentResponse}}
// @Failure 500 {object} dto.BaseResponse
// @Router /posts/{id}/comments [get]
func (c *CommentController) GetAllByPostID(ctx *gin.Context) {
	id := ctx.Param("id")
	page := ctx.Query("page")
	pageSize := ctx.Query("page_size")

	postID, _ := strconv.Atoi(id)
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	queryParams := dto.QueryParams{
		Search:    ctx.Query("search"),
		SortBy:    ctx.Query("sort_by"),
		SortOrder: ctx.Query("sort_order"),
	}

	queryParams.Page = pageInt
	queryParams.PageSize = pageSizeInt

	commentResponses, err := c.commentService.GetAllByPostID(uint(postID), &queryParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve comments"))
		return
	}

	totalCount, err := c.commentService.CountAllByPostID(uint(postID), &queryParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve posts"))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewSuccessResponsePagination(commentResponses, totalCount))
}
