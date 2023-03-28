package config

import "sync"

var oSingle sync.Once

// Instance - Export single instance
func Instance() *Config {
	oSingle.Do(
		func() {
			cfg = new(Config)
			cfg.setName(name)
			cfg.setVersion(version)
		})
	return cfg
}
