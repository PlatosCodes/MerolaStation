package util

import (
	"time"

	"github.com/spf13/viper"
)

type Smtp struct {
	Host     string `mapstructure:"smtpHost"`
	Port     int    `mapstructure:"smtpPort"`
	Username string `mapstructure:"smtpUsername"`
	Password string `mapstructure:"smtpPassword"`
	Sender   string `mapstructure:"smtpSender"`
}

// TODO: Look into changing TokenSymmetricKey to Asymmetric in production.
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	Host                 string        `mapstructure:"smtpHost"`
	Port                 int           `mapstructure:"smtpPort"`
	Username             string        `mapstructure:"smtpUsername"`
	Password             string        `mapstructure:"smtpPassword"`
	Sender               string        `mapstructure:"smtpSender"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
