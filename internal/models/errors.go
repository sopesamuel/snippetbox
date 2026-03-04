package models

import "errors"

var (
	ErrNoRecord = errors.New("Models: No matching record")
	ErrInvalidCredentials = errors.New("Models: Invalid Credentials")
	ErrDuplicateEmail = errors.New("Models: Duplicate Email")
)