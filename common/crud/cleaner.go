package crud

import (
	"github.com/pavlo67/workshop/common/selectors"
)

type Cleaner interface {
	SelectToClean(*Options) (*selectors.Term, error)
	Clean(*selectors.Term, *Options) error
}
