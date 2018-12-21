package identity

import "github.com/pavlo67/punctum/basis/encryption"

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
	Type     CredsType           `json:"type"`
	Cryptype encryption.Cryptype `json:"cryptype"`
	Value    string              `json:"value,omitempty"`
}
