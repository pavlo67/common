package db

import "github.com/pavlo67/common/common/selectors"

type Cleaner interface {
	Clean(term *selectors.Term) error
}
