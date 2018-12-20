package config

type ServerAccess struct {
	Host, User, Pass string
	Port             int
	Path             string // file path, database name, etc,,,
}

// ServerTLS ...
type ServerTLS struct {
	Port        int    `bson:"port"           json:"port"`
	TLSCertFile string `bson:"tls_cert_file"  json:"tls_cert_file"`
	TLSKeyFile  string `bson:"tls_key_file"   json:"tls_key_file"`
}
