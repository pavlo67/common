package demo_server_http_jsschmhr

import (
	"net/http"

	"github.com/pavlo67/partes/fronthttp/componenthtml"
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/identity"
	"github.com/pavlo67/punctum/server_http"
)

var _ server_http.Templator = &templator{}

type templator struct {
	components     map[string]componenthtml.Operator
	htmlLeftStatic string
	htmlNoUser     string
	htmlUserStart  string
	htmlUserEnd    string
	htmlStatic     map[string]string
}

func Templator(components map[string]componenthtml.Operator, joiner program.Joiner) *templator {

	htmlStatic := componenthtml.Static(joiner)

	s := &templator{components: components, htmlStatic: htmlStatic}

	if crudComponentOp, ok := components["crud"]; ok && crudComponentOp != nil {
		s.htmlLeftStatic += "\n" + crudComponentOp.Name() + "<div class=\"ut4\">\n"
		for _, menuItem := range crudComponentOp.Menu("left") {
			s.htmlLeftStatic += "<li><a href=\"" + menuItem.URL + "\">" + menuItem.Label + "</a></li>"
		}

		s.htmlLeftStatic += "</div>\n"
	}

	if confidenterComponentOp, ok := components["confidenter"]; ok && confidenterComponentOp != nil {
		listenersConfidenter := confidenterComponentOp.Listeners()
		endpointsConfidenter := confidenterComponentOp.Endpoints()

		s.htmlNoUser = "<div class=\"ut1\">Хто тут?</div>\n" + `<div class="ut4"><small>` +
			`<input id="login" placeholder="логін" style="width:100px;" /> ` +
			`<input id="password" placeholder="пароль" type="password" style="width:49px;" /> ` +
			`<button id="auth" style="height:20px;width:24px;">&gt;</button> `

		if _, ok := listenersConfidenter["authFB"]; ok {
			s.htmlNoUser += `<a id="authFB"><img src="/images/fb.png" style="margin-top:2px;vertical-align:bottom;"></a> &nbsp; `
		}
		if epBlank, ok := endpointsConfidenter["blank"]; ok {
			s.htmlNoUser += `[<a href="` + epBlank.ServerPath + `">реєстрація</a>]`
		}
		if epForgotPassword, ok := endpointsConfidenter["forgotPassword"]; ok {
			s.htmlNoUser += ` &nbsp;[<a href="` + epForgotPassword.ServerPath + `">забув пароля</a>]`
		}

		s.htmlNoUser += "</small></div>\n"

		if epViewMysself, ok := endpointsConfidenter["viewMyself"]; ok {
			s.htmlUserStart = `<div class="ut1"><b><a href="` + epViewMysself.ServerPath + `">`
			s.htmlUserEnd = `</a></b>`
		} else {
			s.htmlUserStart = `<div class="ut1"><b>`
			s.htmlUserEnd = `</b>`
		}

		// `[<a href="` + epEdit.ServerPath + `">редаґувати профіль</a>] &nbsp; ` +

		if _, ok := endpointsConfidenter["signOut"]; ok {
			s.htmlUserEnd += ` &nbsp; <small>` +
				`[<a id="signOut" href="#">вихід</a>]` +
				`</small></div>`
		}
	}

	return s
}

func (s templator) Context(user *identity.User, _ *http.Request, _ map[string]string) map[string]string {
	var htmlAuth string
	if user == nil || user.ID == "" {
		htmlAuth = s.htmlNoUser
	} else {
		userNickname := user.Nickname
		if userNickname == "" {
			userNickname = "..."
		}
		htmlAuth = s.htmlUserStart + userNickname + s.htmlUserEnd
	}

	context := map[string]string{}
	for k, v := range s.htmlStatic {
		context[k] = v
	}
	context["left"] = s.htmlLeftStatic + htmlAuth

	return context
}
