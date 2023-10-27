package shared

import "github.com/pkg/errors"

var (
	ErrLinkDeleted   = errors.New("link is deleted")
	ErrLinkNotExists = errors.New("link does not exist")
)
