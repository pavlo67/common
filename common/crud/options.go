package crud

import "github.com/pavlo67/workshop/common"

type SaveOptions struct {
	AuthID    common.ID
	Replace   bool
	ReturnIDs bool
}

type GetOptions struct {
	AuthID  common.ID
	GroupBy []string
	OrderBy []string
	Limit0  uint64
	Limit1  uint64
}

type RemoveOptions struct {
	Limit  uint64
	AuthID common.ID
	Delete bool
}
