package store

import "errors"

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrEditConflict       = errors.New("edit conflict")
	ErrInvalidCredentials = errors.New("user store: invalid credentials")
)
