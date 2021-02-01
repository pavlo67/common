package server

// Config ...
type Config struct {
	Port        int    `yaml:"port"          json:"port"`
	NoHTTPS     bool   `yaml:"no_https"      json:"no_https"`
	KeyPath     string `yaml:"key_path"      json:"key_path"`
	TLSCertFile string `yaml:"tls_cert_file" json:"tls_cert_file"`
	TLSKeyFile  string `yaml:"tls_key_file"  json:"tls_key_file"`
}
