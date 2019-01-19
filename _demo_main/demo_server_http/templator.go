package demo_server_http

import (
	"net/http"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/server/server_http"
	"github.com/pavlo67/punctum/starter/joiner"
)

var _ server_http.Templator = &templator{}

type templator struct {
	htmlNoUser string
	htmlMenu   string
}

func newTemplator(joiner joiner.Operator) server_http.Templator {
	htmlMenu := `<linker_server_http><a href="/">root</a></linker_server_http>` + "\n" +
		`<linker_server_http><a href="/section1">section 1</a></linker_server_http>`

	s := &templator{
		htmlNoUser: "user isn't authorized...",
		htmlMenu:   htmlMenu,
	}

	return s
}

func (s templator) Context(user *auth.User, _ *http.Request, _ map[string]string) map[string]string {
	if user != nil && user.ID != "" {
		// return specisic user's template
	}

	context := map[string]string{
		"left": s.htmlNoUser + "\n<p>\n" + s.htmlMenu,
	}

	return context
}
