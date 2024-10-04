package main

import (
	"github.com/dedenfarhanhub/blog-service/internal"
	"github.com/dedenfarhanhub/blog-service/internal/services"
	"github.com/joho/godotenv"
	"log"
)

// @title           Blog API GO
// @version         1.0
// @description     API for blog.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8090
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	// Initialize database connection
	db, err := internal.InitDB()
	if err != nil {
		panic("failed to connect database")
	}

	// Initialize Redis
	redisClient := internal.InitRedis()
	redisService := services.NewRedisService(redisClient)

	r := internal.InitRouter(db, redisService)

	if err := r.Run(":9000"); err != nil {
		log.Fatal(err)
	}
}
