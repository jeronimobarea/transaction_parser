package svcerrors

import "errors"

var (
	ErrConflict   = errors.New("error conflict")
	ErrNotFound   = errors.New("error not found")
	ErrBadRequest = errors.New("error bad request")
)
