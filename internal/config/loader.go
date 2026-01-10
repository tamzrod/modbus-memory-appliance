package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const DefaultConfigPath = "config.yaml"

func LoadOrCreate() (*AppConfig, error) {
	if _, err := os.Stat(DefaultConfigPath); os.IsNotExist(err) {
		if err := writeDefaultConfig(DefaultConfigPath); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(
			"config file not found; default config.yaml created, please review and restart",
		)
	}

	data, err := os.ReadFile(DefaultConfigPath)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.Memory.Validate(); err != nil {
		return nil, err
	}

	if err := cfg.Routing.Validate(cfg.Memory); err != nil {
		return nil, err
	}

	return &cfg, nil
} // ← THIS BRACE MUST EXIST

func writeDefaultConfig(path string) error {
	defaultCfg := AppConfig{
		Memory: MemoryConfig{
			Memories: map[string]MemoryBlock{
				"memory1": {
					Default: true,
					Coils: AreaConfig{
						Start: 0,
						Size:  1024,
					},
					DiscreteInputs: AreaConfig{
						Start: 0,
						Size:  1024,
					},
					HoldingRegisters: AreaConfig{
						Start: 0,
						Size:  4096,
					},
					InputRegisters: AreaConfig{
						Start: 0,
						Size:  4096,
					},
				},
			},
		},
		Routing: RoutingConfig{
			UnitIDMap: map[uint8]string{
				1: "memory1",
			},
		},
	}

	out, err := yaml.Marshal(&defaultCfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, out, 0644)
}
