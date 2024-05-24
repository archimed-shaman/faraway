package config

type NetConfig struct {
	Host          string    `env:"HOST"             yaml:"host"`
	Port          int       `env:"PORT"             yaml:"port"`
	BuffSize      int       `yaml:"buff_size"`
	Timeout       sDuration `yaml:"timeout"`
	MaxConnection int32     `yaml:"max_connections"`
}
