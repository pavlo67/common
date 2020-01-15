package server

// Config ...
type Config struct {
	TLSCertFile string `bson:"tls_cert_file" json:"tls_cert_file"`
	TLSKeyFile  string `bson:"tls_key_file"  json:"tls_key_file"`

	Testers []string `bson:"testers,omitempty" json:"testers,omitempty"`
}
