# Blog API

## Overview

The Blog API is a RESTful service that allows users to create, read, update, and delete blog posts and comments. It provides user authentication and authorization features to secure the endpoints.

## Table of Contents
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Entities](#entities)
- [API Endpoints](#api-endpoints)
- [Code Quality and Organization](#code-quality-and-organization)
- [Features](#features)
- [Security Measures](#security-measures)
- [Creativity and Problem-Solving Approach](#creativity-and-problem-solving-approach)
- [Installation](#installation)
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
 │ └── db.go           # Database initialization
 │ └── redis.go        # Redis initialization
 │ └── migration.go    # Database migration handling
 │ └── routes.go       # Routes definition
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


## Database Tables Structure
### Database Schema
#### users

| Column         | Type         | Index   |
|----------------|--------------|---------|
| `id`           | Integer      | PK      |
| `name`         | String       |         |
| `email`        | String       | Unique  |
| `password_hash`| String       |         |
| `created_at`   | Timestamp    |         |
| `updated_at`   | Timestamp    |         |

#### post

| Column         | Type         | Index   |
|----------------|--------------|---------|
| `id`           | Integer      | PK      |
| `title`        | String       |         |
| `content`      | Text         |         |
| `author_id`    | Integer      | FK(User)|
| `created_at`   | Timestamp    |         |
| `updated_at`   | Timestamp    |         |

#### comments

| Column         | Type         | Index   |
|----------------|--------------|---------|
| `id`           | Integer      | PK      |
| `post_id`      | Integer      | FK(Blog Post)|
| `author_name`  | String       |         |
| `content`      | Text         |         |
| `created_at`   | Timestamp    |         |

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

### Comments
- **POST /posts/{id}/comments** - Add a new comment to a specific post.
- **GET /posts/{id}/comments** - Retrieve all comments associated with a specific post.

### Documentation
- You can access the Swagger documentation at: [Swagger UI](http://localhost:8090/swagger/index.html)

## Code Quality and Organization
- The project follows Go best practices and is compliant with `golint`, ensuring clean and maintainable code.
- It is structured into packages (controllers, services, repositories, etc.) that separate concerns, enhancing readability and testability.

## Features
-  All required features (user registration, authentication, blog CRUD, comments) are complete.

## Security Measures
- **XSS Prevention**: The application prevents Cross-Site Scripting (XSS) attacks by escaping HTML special characters in user-generated content.
- **Security Middleware**: Implemented to enforce security best practices, such as setting security headers.
- **Rate Limiting**: The API incorporates rate limiting to prevent abuse and denial-of-service attacks.
- **HTML Escape**: All special characters are escaped to prevent XSS.

## Creativity and Problem-Solving Approach

In this project, the following creative and problem-solving approaches were applied:

- **Efficient Project Structure**: The project is organized into clean, modular components. Each responsibility, from routing, middleware, database interaction, and services, is separated into distinct directories to promote maintainability and scalability.

- **Post Ownership and Authorization**: To ensure that only the author of a blog post can update or delete it, an ownership check was implemented. This enhances security and ensures proper role-based access control, solving potential issues related to unauthorized modifications.

- **XSS Protection**: All user-generated content, such as blog posts and comments, undergoes HTML sanitization and escaping of special characters. This prevents Cross-Site Scripting (XSS) vulnerabilities and ensures that the API can handle untrusted inputs safely.

- **Rate Limiting**: To protect against brute-force attacks on user login and other endpoints, a rate limiter was implemented. This helps to mitigate potential abuse of the API, ensuring fair usage.

- **Custom Middleware for Security**: Middleware was created to handle authentication using JWT tokens. This ensures that routes which require authentication are protected, while maintaining a clear separation of concerns for security at the application level.

- **Graceful Error Handling**: Thoughtful error handling was applied across API endpoints to provide meaningful error messages to users and maintain consistency in API responses.

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

1. **Create Directory**. Create a directory for your project under `$GOPATH/src/github.com/dedenfarhanhub`:
   ```bash
   mkdir -p $GOPATH/src/github.com/dedenfarhanhub
   ```
2. **Clone the repository**.  Navigate to the project directory. 
   ```bash
   cd $GOPATH/src/github.com/dedenfarhanhub
   git clone https://github.com/dedenfarhanhub/b581007aa20a6533de54b13f08430b5a.git
   ```
3. Rename the cloned repository directory to blog-service
   ```bash
   mv <cloned-repo> blog-service
   ```
4. Create a `.env.docker` file based on the provided sample.
5. Build and run the application using Docker Compose:
   ```bash
   docker-compose up --build
   ```
6. Access the API at `http://localhost:8090`

## License
This project is licensed under the MIT License - see the LICENSE file for details.