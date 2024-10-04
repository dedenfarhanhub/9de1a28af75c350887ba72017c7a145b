package controllers

import (
	"github.com/dedenfarhanhub/blog-service/internal/helpers"
	"net/http"

	"github.com/dedenfarhanhub/blog-service/internal/dto"
	"github.com/dedenfarhanhub/blog-service/internal/services"
	"github.com/gin-gonic/gin"
)

// UserController struct
type UserController struct {
	userService services.UserService
}

// NewUserController controller
func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user with the provided details
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.UserRequest true "User details"
// @Success 200 {object} dto.BaseResponse{data=dto.UserResponse}
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var userDto dto.UserRequest
	if err := ctx.ShouldBindJSON(&userDto); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	userResponse, err := c.userService.Register(&userDto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, userResponse)
}

// Login godoc
// @Summary User login
// @Description Authenticate a user with the provided credentials
// @Tags Users
// @Accept json
// @Produce json
// @Param login body dto.UserLoginRequest true "User login credentials"
// @Success 200 {object} dto.BaseResponse{data=dto.UserResponse}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var loginDto dto.UserLoginRequest
	if err := ctx.ShouldBindJSON(&loginDto); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	userResponse, err := c.userService.Login(&loginDto)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, helpers.NewErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewSuccessResponse(userResponse))
}
