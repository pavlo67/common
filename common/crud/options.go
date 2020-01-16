package crud

import "github.com/pavlo67/workshop/common/identity"

type SaveOptions struct {
	Actor *identity.Key
	// TODO??? check if .Key exists and if it should be existing (insert vs. replace)
}

type GetOptions struct {
	Actor   *identity.Key
	GroupBy []string
	OrderBy []string
	Limit0  uint64
	Limit1  uint64
}

type RemoveOptions struct {
	Actor  *identity.Key
	Limit  uint64
	Delete bool
}
