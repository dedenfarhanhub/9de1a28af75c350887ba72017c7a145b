# Blog API

## Overview

The Blog API is a RESTful service that allows users to create, read, update, and delete blog posts and comments. It provides user authentication and authorization features to secure the endpoints.

## Table of Contents
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Entities](#entities)
- [API Endpoints](#api-endpoints)
- [Security Measures](#security-measures)
- [Code Quality and Organization](#code-quality-and-organization)
- [Running the Project](#running-the-project)
- [License](#license)

## Tech Stack
- **Go**: Programming language used for building the API.
- **Gin**: Web framework for building APIs in Go.
- **GORM**: Object-relational mapping (ORM) library for Go, used for interacting with the MySQL database.
- **MySQL**: Relational database management system for storing user and blog post data.
- **Redis**: In-memory data structure store, used as a caching layer.
- **Docker**: Containerization platform for building and running the application.
- **Swagger**: API documentation and testing tool, providing a user-friendly interface to explore the API endpoints.

## Project Structure
```
blog-service 
├── cmd
 │ ├── server # Main server application entry point 
 │ └── migrate # Migration scripts
├── config # Configuration files
├── db
 │ ├── migrations # Database migration files
├── internal
 │ ├── controllers # API controllers
 │ ├── dto # Data Transfer Objects
 │ └── entities # Database models and ORM definitions
 │ ├── helpers # Helper functions and utilities
 │ ├── middleware # Middleware functions
 │ ├── repositories # Data access layer
 │ ├── services # Business logic and service layer
├── docs # Swagger documentation
├── entrypoint.sh # Entrypoint script for Docker
├── .env.docker # Environment variables for build local
├── .env.docker # Environment variables for Docker
├── Dockerfile # Dockerfile for building the application
└── docker-compose.yml # Docker Compose configuration
└── lint.py # Linting script
```


### Explanation of Project Structure

1. **cmd/**: This directory contains the main application entry points. Each subdirectory represents a different command.
    - **server/**: The entry point for starting the main server application.
    - **migrate/**: Contains scripts for database migration operations.

2. **config/**: This folder holds configuration files that define various settings required for the application, such as environment variables, database connections, and other configurations.

3. **db/**:
    - **migrations/**: Contains migration files used for setting up and altering the database schema.

4. **internal/**: This is where the core business logic of the application resides. It is organized into several subdirectories:
    - **controllers/**: Contains API controllers that handle HTTP requests and responses. Each controller corresponds to a specific resource (e.g., users, posts).
    - **dto/**: Data Transfer Objects that define the structure of request and response payloads for APIs.
    - **entities/**: Defines the database models and ORM (Object-Relational Mapping) entities, representing the structure of the data stored in the database.
    - **helpers/**: Utility functions and helper methods that perform common tasks, such as formatting or data manipulation.
    - **middleware/**: Contains middleware functions for processing requests, such as authentication, logging, or input validation.
    - **repositories/**: The data access layer that interacts with the database, handling CRUD (Create, Read, Update, Delete) operations.
    - **services/**: Implements business logic and interacts with the controllers and repositories to fulfill application requirements.

5. **docs/**: Contains Swagger documentation files for defining and generating API specifications. This helps in documenting the endpoints, request/response structures, and more.

6. **entrypoint.sh**: This is the entrypoint script for the Docker container, which is executed when the container starts. It can be used to run initial setup commands.

7. **.env.docker**: Contains environment variables specific to local development with Docker, such as database credentials and application settings.

8. **Dockerfile**: The file that defines how to build the Docker image for the application, including dependencies and build instructions.

9. **docker-compose.yml**: Defines the services, networks, and volumes for Docker Compose, allowing you to run multiple containers and services with a single command.

10. **lint.py**: A script that helps maintain code quality by checking for linting issues and ensuring adherence to coding standards.

Feel free to adjust any part of the explanation to better fit your project's specifics!


## Database Tables Structure

### Entities

1. **users**
    - **id** (integer, primary key): Unique identifier for the user.
    - **name** (string): Name of the user.
    - **email** (string, unique): Email address of the user.
    - **password_hash** (string): Hashed password for authentication.
    - **created_at** (timestamp): Timestamp of when the user was created.
    - **updated_at** (timestamp): Timestamp of the last update to the user's information.

2. **posts**
- **id** (integer, primary key): Unique identifier for the blog post.
- **title** (string): Title of the blog post.
- **content** (text): Content of the blog post.
- **author_id** (integer, foreign key referencing User): ID of the user who created the post.
- **created_at** (timestamp): Timestamp of when the post was created.
- **updated_at** (timestamp): Timestamp of the last update to the post.
3. **comments**
- **id** (integer, primary key): Unique identifier for the comment.
- **post_id** (integer, foreign key referencing Blog Post): ID of the blog post to which the comment belongs.
- **author_name** (string): Name of the person who wrote the comment.
- **content** (text): Content of the comment.
- **created_at** (timestamp): Timestamp of when the comment was created.

### Indexes Summary
- **Primary Keys**: Each table has a primary key (`id`) that is automatically indexed to ensure unique values and quick access.
- **Unique Index**: The `email` field in the `User` table has a unique index to enforce that no two users can have the same email address.
- **Foreign Keys**: The `author_id` in the `Blog Post` table and `post_id` in the `Comment` table are indexed to improve performance during join operations with the `User` and `Blog Post` tables, respectively.
- **Optional Indexing**: Other columns can be indexed based on application needs, especially if there are frequent search queries on those columns.


## API Endpoints

### User Registration & Authentication
- **POST /register**: Register a new user.
- **POST /login**: Login and receive a token for authentication.

### Blog Posts
- **POST /posts**: Create a new blog post.
- **GET /posts/{id}**: Get blog post details by ID.
- **GET /posts**: List all blog posts.
- **PUT /posts/{id}**: Update a blog post.
- **DELETE /posts/{id}**: Delete a blog post.

### Documentation
- You can access the Swagger documentation at: [Swagger UI](http://localhost:8090/swagger/index.html#/Users/post_register)

## Security Measures
- **XSS Prevention**: The application prevents Cross-Site Scripting (XSS) attacks by escaping HTML special characters in user-generated content.
- **Security Middleware**: Implemented to enforce security best practices, such as setting security headers.
- **Rate Limiting**: The API incorporates rate limiting to prevent abuse and denial-of-service attacks.
- **Authentication**: The service uses token-based authentication (JWT) for securing routes and user management.

## Code Quality and Organization
- The project follows Go best practices and is compliant with `golint`, ensuring clean and maintainable code.
- It is structured into packages (controllers, services, repositories, etc.) that separate concerns, enhancing readability and testability.

## Sample `.env.docker` File

```dotenv
# Database Configuration
DB_USER=root
DB_PASSWORD=mysecretpassword
DB_NAME=mydatabase
DB_HOST=db
DB_PORT=3306

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379

# JWT Configuration
JWT_SECRET=your_jwt_secret_key

# Other Environment Variables
GIN_MODE=release
```
## Installation

1. Clone the repository.
2. Navigate to the project directory. 
3. Create a `.env.docker` file based on the provided sample.
4. Build and run the application using Docker Compose:
```bash
docker-compose up --build
```
5. Access the API at `http://localhost:8090`