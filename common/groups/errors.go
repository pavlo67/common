package groups

import "github.com/pkg/errors"

var ErrNotBelongsTo = errors.New("don't belong to")
var ErrNoRights = errors.New("has no rights")
