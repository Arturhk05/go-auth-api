package errors

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenRevoked = errors.New("refresh token is revoked or invalid")

	ErrAccountInactive = errors.New("account is inactive")
	ErrAccountLocked   = errors.New("account is locked")

	ErrInvalidCredentials = errors.New("invalid email or password")
)
