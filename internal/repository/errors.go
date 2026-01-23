package repository

import "errors"

type RepoError struct{}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserDuplicatedKey = errors.New("user with this email already exists")
	ErrMessageNotFound   = errors.New("message not found")
	ErrInternalServer    = errors.New("internal server error")
	ErrIDRequired        = errors.New("ID is required")
)
