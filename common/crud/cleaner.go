package crud

import (
	"github.com/pavlo67/common/common/selectors"
)

type Cleaner interface {
	Clean(*selectors.Term, *Options) error
}
