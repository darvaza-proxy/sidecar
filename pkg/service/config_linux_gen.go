//go:build linux

package service

//go:generate ./config_linux.sh

// SystemdScript is the custom systemd script.
func (cfg *Config) SystemdScript() string {
	return cfg.GetStringOption("SystemdScript", "")
}

// SetSystemdScript sets the custom systemd script.
func (cfg *Config) SetSystemdScript(script string) {
	cfg.SetOption("SystemdScript", script)
}

// UpstartScript is the custom upstart script.
func (cfg *Config) UpstartScript() string {
	return cfg.GetStringOption("UpstartScript", "")
}

// SetUpstartScript sets the custom upstart script.
func (cfg *Config) SetUpstartScript(script string) {
	cfg.SetOption("UpstartScript", script)
}

// SysvScript is the custom sysv script.
func (cfg *Config) SysvScript() string {
	return cfg.GetStringOption("SysvScript", "")
}

// SetSysvScript sets the custom sysv script.
func (cfg *Config) SetSysvScript(script string) {
	cfg.SetOption("SysvScript", script)
}

// OpenRCScript is the custom openrc script.
func (cfg *Config) OpenRCScript() string {
	return cfg.GetStringOption("OpenRCScript", "")
}

// SetOpenRCScript sets the custom openrc script.
func (cfg *Config) SetOpenRCScript(script string) {
	cfg.SetOption("OpenRCScript", script)
}
