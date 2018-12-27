package identity

import "github.com/pavlo67/punctum/basis/libs/encrlib"

type CredsType string

const CredsToken CredsType = "token"
const CredsPartnerToken CredsType = "partner_token"

const CredsID CredsType = "id"
const CredsNickname CredsType = "nickname"
const CredsEmail CredsType = "email"
const CredsPassword CredsType = "password"
const CredsSentCode CredsType = "sent_code"

const CredsQuestion CredsType = "question"
const CredsQuestionAnswer CredsType = "question_answer"

const CredsAllowedID CredsType = "allowed_id"

const CredsContentToSignature CredsType = "content_to_signature"
const CredsSignature CredsType = "signature"
const CredsPublicKeyAddress CredsType = "public_key_address"
const CredsPublicKey CredsType = "public_key"
const CredsPrivateKey CredsType = "private_key"

type Creds struct {
	Type     CredsType        `json:"type"`
	Cryptype encrlib.Cryptype `json:"cryptype"`
	Value    string           `json:"value,omitempty"`
}
