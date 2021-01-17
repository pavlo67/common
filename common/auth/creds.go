package auth

import (
	"github.com/pavlo67/workshop/common"
)

type CredsType = string

const CredsToSet CredsType = "to_set"
const CredsRealm CredsType = "realm"

const CredsIP CredsType = "ip"

const CredsJWT CredsType = "jwt"
const CredsJWTRefresh CredsType = "jwt_refresh"

const CredsToken CredsType = "token"
const CredsPartnerToken CredsType = "partner_token"

const CredsRoles CredsType = "roles"

const CredsLogin CredsType = "login"
const CredsNickname CredsType = "nickname"
const CredsEmail CredsType = "email"
const CredsPassword CredsType = "password"
const CredsTemporaryKey CredsType = "temporary_key"

const CredsQuestion CredsType = "question"
const CredsQuestionAnswer CredsType = "question_answer"

const CredsAllowedID CredsType = "allowed_id"

const CredsKeyToSignature CredsType = "key_to_signature"
const CredsSignature CredsType = "signature"
const CredsPublicKeyBase58 CredsType = "public_key_base58"
const CredsPublicKeyEncoding CredsType = "public_key_encoding"
const CredsPrivateKey CredsType = "private_key"

const CredsCompanyID CredsType = "company_id"
const CredsCompanyIDExternal CredsType = "company_id_external"

const CredsPasshash CredsType = "passhash"
const CredsPasshashCryptype CredsType = "passhash_cryptype"

type Creds = common.Map
