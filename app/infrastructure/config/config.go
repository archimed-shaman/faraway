package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type Config struct {
	App AppConfig `yaml:"app"`
	Log LogConfig `yaml:"log"`
	Net NetConfig `yaml:"net"`
}

var configInst *Config

func NewConfig(path string, logger *zap.Logger) *Config {
	configInst = new(Config)

	logger.Info("Reading config", zap.String("path", path))

	err := cleanenv.ReadConfig(path, configInst)
	if err != nil {
		panic(err)
	}

	logConfig(configInst, logger)

	return configInst
}

func Get() *Config {
	if configInst == nil {
		panic("config is not initialized")
	}

	return configInst
}

func logConfig(cfg *Config, logger *zap.Logger) {
	val := reflect.ValueOf(cfg).Elem()
	logStruct("", val, logger)
}

func logStruct(prefix string, val reflect.Value, logger *zap.Logger) {
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)

		typeField := val.Type().Field(i)
		tag := typeField.Tag
		name := typeField.Name

		displayedName := tag.Get("yaml")
		if displayedName == "" {
			displayedName = name
		}

		if valueField.Kind() == reflect.Struct {
			logStruct(prefix+":"+displayedName, valueField, logger)
		} else {
			displayedValue := fmt.Sprintf("%v", valueField)

			if isSecret(displayedName) && displayedValue != "" {
				displayedValue = displaySecret(displayedValue)
			}

			logger.Info("Config", zap.String(prefix, displayedValue))
		}
	}
}

func isSecret(name string) bool {
	if strings.Contains(strings.ToLower(name), "key") {
		return true
	}

	if strings.Contains(strings.ToLower(name), "salt") {
		return true
	}

	if strings.Contains(strings.ToLower(name), "password") {
		return true
	}

	return false
}

func displaySecret(secret string) string {
	return strings.Repeat("*", len(secret))
}
