package dataimporter

import (
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/identity"
)

const ActionKey crud.ActionKey = "imported"

func Identity(sourceURL, originalID string) *identity.Item {
	ident := identity.FromURLRaw(sourceURL)
	ident.ID = originalID

	return &ident
}

func SourceKey(history crud.History) string { //  , sourceTime time.Time
	for _, action := range history {
		if action.Key == ActionKey && action.Identity != nil {
			return action.Identity.Domain + "/" + action.Identity.Path + "##" + action.Identity.ID // , action.DoneAt
		}
	}

	return ""
}
