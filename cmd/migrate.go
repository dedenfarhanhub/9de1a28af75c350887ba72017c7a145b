package main

import (
	"github.com/dedenfarhanhub/blog-service/internal"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize database connection
	db, err := internal.InitDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Handle command-line arguments to determine which migration command to run
	args := os.Args
	if len(args) < 2 {
		log.Fatal("please provide a migration command: 'up' or 'down'")
	}

	command := args[1]

	switch command {
	case "up":
		// Run migrations
		if err := internal.RunMigrations(db); err != nil {
			log.Fatalf("could not run migrations: %v", err)
		}
		log.Println("Migrations applied successfully")
	case "down":
		// Rollback migrations
		if len(args) < 3 {
			log.Fatal("please provide the number of steps for rolling back")
		}

		steps := args[2]

		if err := internal.RollbackMigrations(db, steps); err != nil {
			log.Fatalf("could not roll back migrations: %v", err)
		}
		log.Println("Migrations rolled back successfully")
	default:
		log.Fatalf("unknown migration command: %s", command)
	}
}
