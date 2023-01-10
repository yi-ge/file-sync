package config

// Config - Export Base Config struct
type Config struct {
	name    string
	version string
}

// SystemConfig - Export System Config struct
type SystemConfig struct {
	LogPath  string
	ExePath  string
	WorkPath string
}

var cfg *Config

// SConfig - Export System Config
var SConfig SystemConfig

// setName - Set the name
func (p *Config) setName(name string) {
	p.name = name
}

// GetName - Get the name
func (p *Config) GetName() string {
	return p.name
}

// setVersion - Set the version
func (p *Config) setVersion(version string) {
	p.version = version
}

// GetVersion - Get the version
func (p *Config) GetVersion() string {
	return p.version
}
