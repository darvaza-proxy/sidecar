package certmagic

// Config ...
type Config struct {
	Key    string `default:"key.pem"`
	Cert   string `default:"cert.pem"`
	CARoot string `default:"caroot.pem"`
}

// New creates a new [storage.Store] using the [Config] values
func (cfg *Config) New(options ...OptionFunc) (*Store, error) {
	var opts []OptionFunc

	if cfg.CARoot != "" {
		opts = append(opts, WithTrustedRoots(cfg.CARoot))
	}

	if cfg.Key != "" {
		var certs []string
		if cfg.Cert != "" {
			certs = append(certs, cfg.Cert)
		}

		opts = append(opts, WithKey(cfg.Key, certs...))
	}

	opts = append(opts, options...)
	return New(opts...)
}
