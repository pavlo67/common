package auth

import "github.com/pavlo67/workshop/common/libraries/encrlib"

type CredsType string

const CredsIP CredsType = "ip"

const CredsJWT CredsType = "jwt"

const CredsToken CredsType = "token"
const CredsPartnerToken CredsType = "partner_token"

const CredsIdentityKey CredsType = "identity_key"
const CredsLogin CredsType = "login"
const CredsNickname CredsType = "nickname"
const CredsEmail CredsType = "email"
const CredsPassword CredsType = "password"
const CredsSentCode CredsType = "sent_code"

const CredsQuestion CredsType = "question"
const CredsQuestionAnswer CredsType = "question_answer"

const CredsAllowedID CredsType = "allowed_id"

const CredsKeyToSignature CredsType = "key_to_signature"
const CredsSignature CredsType = "signature"
const CredsPublicKeyBase58 CredsType = "public_key_base58"
const CredsPublicKeyEncoding CredsType = "public_key_encoding"
const CredsPrivateKey CredsType = "private_key"

type Values map[CredsType]string

type Creds struct {
	Cryptype encrlib.Cryptype `json:",omitempty"`
	Values   Values           `json:",omitempty"`
}
