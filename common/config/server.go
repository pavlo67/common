package config

type ServerAccess struct {
	Host string `bson:"host,omitempty" json:"host,omitempty"`
	User string `bson:"user,omitempty" json:"user,omitempty"`
	Pass string `bson:"pass,omitempty" json:"pass,omitempty"`
	Port int    `bson:"port,omitempty" json:"port,omitempty"`
	Path string `bson:"path,omitempty" json:"path,omitempty"` // file path, database name, etc,,,
}

// Server ...
type Server struct {
	Port int `bson:"port" json:"port"`

	TLSCertFile string `bson:"tls_cert_file" json:"tls_cert_file"`
	TLSKeyFile  string `bson:"tls_key_file"  json:"tls_key_file"`

	Testers []string `bson:"testers,omitempty" json:"testers,omitempty"`
}
