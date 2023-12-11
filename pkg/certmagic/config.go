package certmagic

// Config ...
type Config struct {
	Key    string `default:"key.pem"`
	Cert   string `default:"cert.pem"`
	CARoot string `default:"caroot.pem"`

	DefaultIssuer  string `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	FallbackDomain string `json:",omitempty" yaml:",omitempty" toml:",omitempty"`

	Accounts []AccountConfig `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
}

// AccountConfig ...
type AccountConfig struct {
	Name   string
	EMail  string
	Issuer string

	Providers []ProviderConfig
}

// ProviderConfig ...
type ProviderConfig struct {
	Name  string
	Token string

	Domains map[string]*DomainConfig
}

// DomainConfig ...
type DomainConfig struct {
	Name string
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
