package rherr

import "errors"

var (
	ErrBadRequest    = errors.New("bad request")
	ErrForbidden     = errors.New("forbidden")
	ErrNotFound      = errors.New("not found")
	ErrInternal      = errors.New("internal server error")
	ErrProfileExists = errors.New("profile exists")
	ErrUsernameTaken = errors.New("username in use")
)
