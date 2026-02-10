package config

import "github.com/kelseyhightower/envconfig"

const envPrefix = "app"

type DB struct {
	Host     string `required:"true"`
	User     string `required:"true"`
	Password string `required:"true"`
	Name     string `required:"true"`
	SSLMode  string `required:"true"`
	Port     int    `required:"true" default:"5432"`
}

type Config struct {
	DB DB
}

func BuildConfig() (*Config, error) {
	c := &Config{}

	err := envconfig.Process(envPrefix, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
