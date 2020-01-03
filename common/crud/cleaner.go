package crud

import "github.com/pavlo67/workshop/common/selectors"

type Cleaner interface {
	SelectToClean(*RemoveOptions) (*selectors.Term, error)
	Clean(*selectors.Term, *RemoveOptions) error
}
