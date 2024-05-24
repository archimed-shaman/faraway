package config

type LogConfig struct {
	Level string `env:"LOG_LEVEL" yaml:"level"`
}
