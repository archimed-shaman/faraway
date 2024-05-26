package config

type AppConfig struct {
	Window               sDuration `yaml:"window"`
	MaxDifficulty        int       `yaml:"max_difficulty"`
	RateDifficultyFactor float64   `yaml:"rate_difficulty_factor"`
}
