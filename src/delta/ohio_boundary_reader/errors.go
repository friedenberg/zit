package ohio_boundary_reader

import (
	"errors"
)

var (
	ErrBoundaryNotFound      = errors.New("boundary not found")
	ErrExpectedContentRead   = errors.New("expected content read")
	ErrExpectedBoundaryRead  = errors.New("expected boundary read")
	ErrReadFromSmallOverflow = errors.New(
		"reader provided more bytes than max int",
	)
	ErrInvalidBoundaryReaderState = errors.New("invalid boundary reader state")
)
