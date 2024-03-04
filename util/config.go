package util

import (
	"time"

	"github.com/spf13/viper"
)

// TODO: Look into changing TokenSymmetricKey to Asymmetric in production.
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	FrontendAddress3     string        `mapstructure:"FRONTEND_ADDRESS_3"`
	FrontendAddress4     string        `mapstructure:"FRONTEND_ADDRESS_4"`
	SmtpHost             string        `mapstructure:"SMTP_HOST"`
	SmtpPort             int           `mapstructure:"SMTP_PORT"`
	SmtpUsername         string        `mapstructure:"SMTP_USERNAME"`
	SmtpPassword         string        `mapstructure:"SMTP_PASSWORD"`
	SmtpSender           string        `mapstructure:"SMTP_SENDER"`
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
