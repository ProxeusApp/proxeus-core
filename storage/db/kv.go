package db

import (
	"errors"
)

var (
	ErrNotFound = errors.New("key does not exist")
)
