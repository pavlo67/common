package crud

import (
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/selectors"
)

type Options struct {
	Identity *auth.Identity
	Selector selectors.Item
}

func (options *Options) GetIdentity() *auth.Identity {
	if options == nil {
		return nil
	}
	return options.Identity
}

func (options *Options) WithSelector(selector selectors.Item) *Options {
	if options == nil {
		return &Options{Selector: selector}
	}
	optionsCopied := *options
	optionsCopied.Selector = selector

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
