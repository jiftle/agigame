package server

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// AgentConfig holds the agent layer settings. Consumed by the agent package
// from P2 onwards; the websocket hub forwards them so the UI can select a mode.
type AgentConfig struct {
	Mode          string `yaml:"mode"`             // rules | llm | hybrid | manual
	LLMIntervalMS int    `yaml:"llm_interval_ms"`  // ms between LLM decisions
	SafetyNetFrames int `yaml:"safety_net_frames"` // frames of no progress before rescue (0 = 120)
	Model         string `yaml:"model"`
	APIKeyEnv     string `yaml:"api_key_env"`
	RulesFirst    bool   `yaml:"rules_first"`
}

// Config is the root configuration for the GoBoy-LLM server.
type Config struct {
	ROM   string            `yaml:"rom"`
	Server ServerConfig     `yaml:"server"`
	Emulator EmulatorConfig `yaml:"emulator"`
	Agent AgentConfig       `yaml:"agent"`
	Game  string            `yaml:"game"`     // game plugin id, e.g. "sml"
	WebUI string            `yaml:"webui"`    // webui asset directory
}

// ServerConfig holds the HTTP/WebSocket server settings.
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// EmulatorConfig holds emulator runtime settings.
type EmulatorConfig struct {
	FPS       int    `yaml:"fps"`        // emulation speed (60)
	FrameSkip int    `yaml:"frame_skip"` // N frames between webui frame pushes
	SavePath  string `yaml:"save_path"`
	CGB       bool   `yaml:"cgb"`
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		ROM:   "./roms/super_mario_land.gb",
		WebUI: "./webui",
		Game:  "sml",
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Emulator: EmulatorConfig{
			FPS:       60,
			FrameSkip: 4,
			SavePath:  "./saves/",
		},
		Agent: AgentConfig{
			Mode:             "manual",
			LLMIntervalMS:    500,
			SafetyNetFrames:  120,
			Model:            "gpt-4o-mini",
			APIKeyEnv:        "OPENAI_API_KEY",
			RulesFirst:       true,
		},
	}
}

// LoadConfig reads a YAML config file and merges it over the defaults.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}