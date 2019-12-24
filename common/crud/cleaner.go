package crud

import "github.com/pavlo67/workshop/common/selectors"

type Cleaner interface {
	Clean(*selectors.Term, *RemoveOptions) error
}
