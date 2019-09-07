package a

import (
	"github.com/pavlo67/punctum/basis/encrlib"
	"github.com/pavlo67/punctum/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "auth"

type Message struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}

type Callback string

const Confirm Callback = "confirm"
const SendCode Callback = "send_code"

type Operator interface {
	// Create stores registration data and (as usual) sends confirmation code to user.
	Create(creds ...Creds) ([]Message, error)

	// Use can require multi-steps...
	// It autenticates the user (using his own or admin's creds) and can also add, update or remove user's creds.
	// "User" are removed at all if no its creds remains.
	Use(toUse, toAuth Creds, toSet ...Creds) (*User, []Message, error)

	AddCallback(key Callback, url string)
}

// Creds --------------------------------------------------------------------------------

//type Schema string
//const Password Schema = "password"
//const Session Schema = "session"
//const Partner Schema = "partner"
//const SentCode Schema = "sent_code"
//const Question Schema = "question"

type CredsType string

const CredsID CredsType = "id"
const CredsNickname CredsType = "nickname"
const CredsEmail CredsType = "email"
const CredsPassword CredsType = "password"
const CredsToken CredsType = "token"
const CredsPartnerToken CredsType = "partner"
const CredsSentCode CredsType = "code"
const CredsQuestionAnswer CredsType = "question"

type Creds struct {
	Type     CredsType        `json:"type"`
	Cryptype encrlib.Cryptype `json:"cryptype"`
	Values   []string         `json:"values,omitempty"`
}

func (c *Creds) FirstValue() string {
	if c != nil && len(c.Values) >= 1 {
		return c.Values[0]
	}

	return ""
}

func (c *Creds) SecondValue() string {
	if c != nil && len(c.Values) >= 2 {
		return c.Values[1]
	}

	return ""
}

// User -----------------------------------------------------------------------------------------

const UserPath = "user"

type User struct {
	ID       string `bson:"id"       json:"id"`
	Nickname string `bson:"nickname" json:"nickname"`
}

func (user *User) Identity() *auth.IDentity {
	if user == nil {
		return nil
	}
	return &auth.IDentity{
		Domain: joiner.SystemDomain(),
		Path:   UserPath,
		ID:     user.ID,
	}
}
