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
- [ ] User login with email and password
- [ ] JWT token generation and validation
- [ ] Password hashing and verification
- [ ] CORS support

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
