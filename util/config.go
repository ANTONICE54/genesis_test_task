package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver        string `mapstructure:"DB_DRIVER"`
	DBSource        string `mapstructure:"DB_SOURCE"`
	MigrationURL    string `mapstructure:"MIGRATION_URL"`
	ServerAddress   string `mapstructure:"SERVER_ADDRESS"`
	MailerDomain    string `mapstructure:"MAILER_DOMAIN"`
	MailerHost      string `mapstructure:"MAILER_HOST"`
	MailerPort      int    `mapstructure:"MAILER_PORT"`
	DailyEmailsTime string `mapstructure:"DAILY_EMAILS_TIME"`
	RateAPIKey      string `mapstructure:"RATE_API_KEY"`
}

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
