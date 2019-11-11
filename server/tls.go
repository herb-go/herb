package server

//TLSConfig tls config
type TLSConfig struct {
	//TLS whether use tls
	TLS bool
	//TLSCertPath tls cert file path
	TLSCertPath string
	//TLSKeyPath tls key file path
	TLSKeyPath string
}

//ServerIsTLSEnabeld return is server tls enabled
func (c *TLSConfig) ServerIsTLSEnabeld() bool {
	return c.TLS
}
