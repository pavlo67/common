package crud

import (
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/selectors"
	"github.com/pavlo67/common/common/selectors/logic"
)

type Options struct {
	Identity *auth.Identity
	Selector *selectors.Term
	Ranges   *Ranges
}

type Ranges struct {
	GroupBy []string
	OrderBy []string
	JoinTo  string
	Values  []interface{}
	Offset  uint64
	Limit   uint64
}

func (options *Options) GetIdentity() *auth.Identity {
	if options == nil {
		return nil
	}
	return options.Identity
}

func (options *Options) WithSelector(selector *selectors.Term) *Options {
	if options == nil {
		return &Options{Selector: selector}
	}
	optionsCopied := *options

	if options.Selector == nil {
		optionsCopied.Selector = selector
	} else if selector != nil {
		optionsCopied.Selector = logic.AND(selector, options.Selector)
	}

	return &optionsCopied
}

func (options *Options) WithRanges(Ranges *Ranges) *Options {
	if options == nil {
		return &Options{Ranges: Ranges}
	}
	optionsCopied := *options
	options.Ranges = Ranges

	return &optionsCopied
}

func (options *Options) HasRole(oneOfRoles ...rbac.Role) bool {
	if options == nil || options.Identity == nil {
		return false
	}

	return options.Identity.Roles.Has(oneOfRoles...)
}

func OptionsWithRoles(roles ...rbac.Role) *Options {
	return &Options{
		Identity: &auth.Identity{
			Roles: roles,
		},
	}
}
