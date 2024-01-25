//go:build linux

package service

// UserService indicates the service should install as current user.
func (cfg *Config) UserService() bool {
	return cfg.GetBoolOption("UserService", false)
}

// SetUserService sets the service to install as current user.
func (cfg *Config) SetUserService(value bool) {
	cfg.SetOption("UserService", value)
}

func (*Config) setOSDefaults() {}
