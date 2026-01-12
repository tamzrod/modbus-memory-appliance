// internal/config/ipfilter.go


package config

type IPFilterConfig struct {
	Allow []string `yaml:"allow"`
	Deny  []string `yaml:"deny"`
}
