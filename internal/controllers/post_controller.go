package controllers

import (
	"github.com/dedenfarhanhub/blog-service/internal/dto"
	"github.com/dedenfarhanhub/blog-service/internal/helpers"
	"github.com/dedenfarhanhub/blog-service/internal/services"
	"github.com/gin-gonic/gin"
	"html"
	"net/http"
	"strconv"
)

// PostController struct
type PostController struct {
	postService services.PostService
}

// NewPostController controller
func NewPostController(postService services.PostService) *PostController {
	return &PostController{postService: postService}
}

// Create godoc
// @Summary Create a new post
// @Description Create a new post with the given details
// @Tags Posts
// @Accept json
// @Produce json
// @Param postRequest body dto.PostRequest true "Post Request"
// @Success 200 {object} dto.BaseResponse{data=dto.PostResponse}
// @Failure 400 {object} dto.BaseResponse
// @Router /posts [post]
// @Security BearerAuth
func (c *PostController) Create(ctx *gin.Context) {
	var postRequest dto.PostRequest
	if err := ctx.ShouldBindJSON(&postRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(400, "Invalid request payload"))
		return
	}

	userID := ctx.MustGet("userID").(uint)
	postRequest.AuthorID = userID
	// Sanitasi input untuk mencegah XSS
	postRequest.Title = html.EscapeString(postRequest.Title)
	postRequest.Content = html.EscapeString(postRequest.Content)
	postResponse, err := c.postService.CreatePost(&postRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(400, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewSuccessResponse(postResponse))
}

// Update godoc
// @Summary Update a post
// @Description Update an existing post
// @Tags Posts
// @Accept  json
// @Produce  json
// @Param id path int true "Post ID"
// @Param postRequest body dto.PostRequest true "Post Request"
// @Success 200 {object} dto.BaseResponse{data=dto.PostResponse}
// @Failure 400 {object} dto.BaseResponse
// @Router /posts/{id} [put]
// @Security BearerAuth
func (c *PostController) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(400, "Invalid post ID"))
		return
	}

	var postRequest dto.PostRequest
	if err := ctx.ShouldBindJSON(&postRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(400, "Invalid request payload"))
		return
	}

	userID := ctx.MustGet("userID").(uint)
	postRequest.AuthorID = userID
	// Sanitasi input untuk mencegah XSS
	postRequest.Title = html.EscapeString(postRequest.Title)
	postRequest.Content = html.EscapeString(postRequest.Content)
	postResponse, err := c.postService.Update(uint(id), &postRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(400, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewSuccessResponse(postResponse))
}

// GetByID godoc
// @Summary Get a post by ID
// @Description Get details of a post by its ID
// @Tags Posts
// @Produce  json
// @Param id path int true "Post ID"
// @Success 200 {object} dto.BaseResponse{data=dto.PostResponse}
// @Failure 400 {object} dto.BaseResponse
// @Router /posts/{id} [get]
func (c *PostController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(400, "Invalid post ID"))
		return
	}

	postResponse, err := c.postService.GetPostByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, helpers.NewErrorResponse(404, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewSuccessResponse(postResponse))
}

// Delete godoc
// @Summary Delete a post by ID
// @Description Delete a post by its ID
// @Tags Posts
// @Produce  json
// @Param id path int true "Post ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Router /posts/{id} [delete]
// @Security BearerAuth
func (c *PostController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(400, "Invalid post ID"))
		return
	}

	userID := ctx.MustGet("userID").(uint)
	err = c.postService.Delete(uint(id), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, helpers.NewErrorResponse(404, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewSuccessResponse(nil))
}

// GetAll godoc
// @Summary Get all posts
// @Description Retrieve all posts with filters
// @Tags Posts
// @Produce  json
// @Param search query string false "Search by title or content"
// @Param sort_by query string false "Sort by field"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.BaseResponse{data=dto.PaginationResponse{items=[]dto.PostResponse}}
// @Failure 500 {object} dto.BaseResponse
// @Router /posts [get]
func (c *PostController) GetAll(ctx *gin.Context) {
	queryParams := dto.QueryParams{
		Search:    ctx.Query("search"),
		SortBy:    ctx.Query("sort_by"),
		SortOrder: ctx.Query("sort_order"),
	}

	// Optionally, convert page and pageSize to integers
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	queryParams.Page = page
	queryParams.PageSize = pageSize

	postResponses, err := c.postService.GetAll(&queryParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve posts"))
		return
	}

	totalCount, err := c.postService.Count(&queryParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve posts"))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewSuccessResponsePagination(postResponses, totalCount))
}
