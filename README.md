# AUTH API with Golang

This is an project built with Golang that provides an authentication API implementing good practices.

> This project is intended for study purposes and feedback is welcome. If you have any suggestions for improvements, please feel free to open an issue or submit a pull request.

## TDD

In this project I tried to follow Test-Driven Development (TDD) principles, writing tests before implementing the actual functionality. This approach helps ensure that the code is well-tested and meets the specified requirements.

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

- [ ] User registration with email and password
- [ ] User login with email and password
- [ ] JWT token generation and validation
- [ ] Password hashing and verification
- [ ] CORS support