package internal

import (
	"github.com/dedenfarhanhub/blog-service/docs"
	"github.com/dedenfarhanhub/blog-service/internal/controllers"
	"github.com/dedenfarhanhub/blog-service/internal/middleware"
	"github.com/dedenfarhanhub/blog-service/internal/repositories"
	"github.com/dedenfarhanhub/blog-service/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
	"time"
)

// InitRouter initializes the Gin router with routes and middleware.
func InitRouter(db *gorm.DB, redisService *services.RedisService) *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	if gin.Mode() == gin.ReleaseMode {
		r.Use(middleware.SecurityMiddleware())
		r.Use(middleware.XSS())
	}
	r.Use(middleware.Cors())
	r.Use(middleware.RateLimiter(rate.Every(1*time.Minute), 60)) // 60 requests per minute

	// Swagger endpoint
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	postRepo := repositories.NewPostRepository(db)
	commentRepo := repositories.NewCommentRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, redisService)
	postService := services.NewPostService(postRepo, userService, redisService)
	commentService := services.NewCommentService(commentRepo, postService)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	postController := controllers.NewPostController(postService)
	commentController := controllers.NewCommentController(commentService)

	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)

	// Post Routes
	postGroup := r.Group("/posts")
	{
		postGroup.POST("/", middleware.AuthMiddleware(), postController.Create)
		postGroup.GET("/:id", postController.GetByID)
		postGroup.GET("/", postController.GetAll)
		postGroup.PUT("/:id", middleware.AuthMiddleware(), postController.Update)
		postGroup.DELETE("/:id", middleware.AuthMiddleware(), postController.Delete)

		// Comment Routes nested under Post
		postGroup.POST("/:id/comments", commentController.Create)
		postGroup.GET("/:id/comments", commentController.GetAllByPostID)
	}

	return r
}
