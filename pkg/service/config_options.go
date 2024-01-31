package service

import (
	"github.com/kardianos/service"
)

// GetOption gets an arbitrary option from the [service.Config]'s
// KeyValue store.
func (cfg *Config) GetOption(name string) (any, bool) {
	if cfg.Option == nil {
		return nil, false
	}

	v, found := cfg.Config.Option[name]
	return v, found
}

// SetOption sets an arbitrary option on the [service.Config]'s
// KeyValue store.
func (cfg *Config) SetOption(name string, value any) {
	if cfg.Option == nil {
		cfg.Config.Option = make(service.KeyValue)
	}

	cfg.Option[name] = value
}

// GetBoolOption gets a boolean option from the [service.Config]'s
// KeyValue store.
func (cfg *Config) GetBoolOption(name string, defaultValue bool) bool {
	if vi, found := cfg.GetOption(name); found {
		if v, is := vi.(bool); is {
			return v
		}
	}

	return defaultValue
}

// SetBoolOption sets a boolean option to the [service.Config]'s
// KeyValue store.
func (cfg *Config) SetBoolOption(name string, value bool) {
	cfg.SetOption(name, value)
}

// GetStringOption gets a string option from the [service.Config]'s
// KeyValue store.
func (cfg *Config) GetStringOption(name, defaultValue string) string {
	if vi, found := cfg.GetOption(name); found {
		if v, is := vi.(string); is {
			return v
		}
	}

	return defaultValue
}

// SetStringOption sets a string option to the [service.Config]'s
// KeyValue store.
func (cfg *Config) SetStringOption(name, value string) {
	cfg.SetOption(name, value)
}

// GetEnv gets an environment variable from the [service.Config].
func (cfg *Config) GetEnv(name, defaultValue string) string {
	if cfg.EnvVars != nil {
		if s, ok := cfg.EnvVars[name]; ok {
			return s
		}
	}

	return defaultValue
}

// SetEnv sets an environment variable on the [service.Config].
func (cfg *Config) SetEnv(name, value string) {
	if cfg.EnvVars == nil {
		cfg.EnvVars = make(map[string]string)
	}

	cfg.EnvVars[name] = value
}
