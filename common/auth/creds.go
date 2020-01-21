package auth

type CredsType string

const CredsToSet CredsType = "to_set"

const CredsIP CredsType = "ip"

const CredsJWT CredsType = "jwt"

const CredsToken CredsType = "token"
const CredsPartnerToken CredsType = "partner_token"

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

const CredsPasshash CredsType = "passhash"
const CredsPasshashCryptype CredsType = "passhash_cryptype"

type Creds map[CredsType]string
