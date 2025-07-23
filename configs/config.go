package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int
	}
	Notifier struct {
		UseTelegram     bool
		UseSlack        bool
		TelegramToken   string
		TelegramChatID  string
		SlackWebhookURL string
	}
	Git struct {
		CloneBaseDir string
	}
	Docker struct {
		Network string
	}
	Database struct {
		Driver string
		DSN    string
	}
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv() 

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
