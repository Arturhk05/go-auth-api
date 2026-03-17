# AUTH API with Golang

This is an project built with Golang that provides an authentication API implementing good practices.

### Packages

- **github.com/gin-gonic/gin** - Web framework
- **github.com/lib/pq** - PostgreSQL driver
- **github.com/golang-jwt/jwt/v5** - JWT authentication
- **golang.org/x/crypto/bcrypt** - Password hashing
- **github.com/joho/godotenv** - Environment variables
- **github.com/go-playground/validator/v10** - Data validation
- **github.com/google/uuid** - UUID generation
- **github.com/gin-contrib/cors** - CORS middleware

## Features

- [x] User registration with email and password
- [x] User login with email and password
- [x] JWT token generation and validation
- [x] Password hashing and verification
- [x] CORS support
- [x] CI/CD pipeline setup
- [x] Implement refresh tokens
- [x] Add Swagger documentation
- [ ] Improve error handling and response formatting
- [ ] Add google OAuth authentication
- [ ] Add rate limiting
- [ ] Add email verification
- [ ] Add password reset functionality
- [ ] Add user roles and permissions
- [ ] Add logging and monitoring
- [ ] Redis caching for session management
- [ ] Add unit and integration tests for all components


## Project Structure

```
AuthGo/
├── .vscode/                      # VS Code workspace settings
├── cmd/                          # Command-line applications
│   ├── api/                      # Main API server
│   │   └── main.go              # Entry point for the API
│   └── migrate/                  # Database migration CLI
│       └── main.go              # Migration runner
├── config/                       # Configuration management
│   ├── config.go                # Configuration loading and setup
│   └── config_test.go           # Config tests
├── database/                     # Database layer
│   ├── database.go              # Database connection
│   └── migrations/               # SQL migration files
├── internal/                     # Private application code
│   ├── handlers/                 # HTTP request handlers
│   ├── middlewares/              # HTTP middlewares (auth, CORS, etc)
│   ├── models/                   # Data models/entities
│   ├── repositories/             # Data access layer
│   ├── services/                 # Business logic
│   └── utils/                    # Utility functions
├── go.mod                       # Go module definition
├── go.sum                       # Go dependencies checksums
└── README.md                    # This file
```

## Notes

- I have roll back my code to a previous commit, because I had no idea about what my code was doing, and I wanted to start over with a better understanding of the project structure and the code itself. (12/06/2024)
- I give up on TDD, I dont have experience with Golang and I started to struggle with writing tests before writing the actual code, I will write the code first and then write the tests later. (13/06/2024)

## Setup Development

1. Create a `.env` file in the root directory:

```bash
copy .env.example .env
```

2. Run docker compose to start the PostgreSQL database:

```bash
docker compose up -d
```

3. Run migrations to create the database schema:

```bash
make migrate-up
```

4. Start the API server:

```bash
make run-api
# or
go run cmd/api/main.go
```
