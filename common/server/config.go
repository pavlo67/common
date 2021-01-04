package server

// Config ...
type Config struct {
	Port        int      `yaml:"port"          json:"port"`
	TLSCertFile string   `yaml:"tls_cert_file" json:"tls_cert_file"`
	TLSKeyFile  string   `yaml:"tls_key_file"  json:"tls_key_file"`
	Testers     []string `yaml:"testers"       json:"testers,omitempty"`
}
