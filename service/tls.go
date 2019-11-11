package service

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

//ServerTLSCertPath return server tls cert path
func (c *TLSConfig) ServerTLSCertPath() string {
	return c.TLSCertPath
}

//ServerTLSKeyPath resturn serve tls key path
func (c *TLSConfig) ServerTLSKeyPath() string {
	return c.TLSKeyPath
}
