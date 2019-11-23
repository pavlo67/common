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
}

type RemoveOptions struct {
	AuthID common.ID
	Delete bool
}
