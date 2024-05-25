package config

type AppConfig struct {
	MaxDifficulty        int     `yaml:"max_difficulty"`
	RateDifficultyFactor float64 `yaml:"rate_difficulty_factor"`
}
