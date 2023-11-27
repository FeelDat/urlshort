package shared

import "github.com/pkg/errors"

// ErrLinkDeleted is returned when a link is deleted.
// ErrLinkNotExists is returned when a link does not exist.
var (
	ErrLinkDeleted   = errors.New("link is deleted")
	ErrLinkNotExists = errors.New("link does not exist")
)
