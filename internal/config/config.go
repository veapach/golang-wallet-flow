package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DBHost   string `env:"DB_HOST" envDefault:"localhost"`
	DBPort   int    `env:"DB_PORT" envDefault:"5432"`
	DBUser   string `env:"DB_USER,required"`
	DBPass   string `env:"DB_PASS,required"`
	DBName   string `env:"DB_NAME,required"`
	HTTPAddr string `env:"HTTP_ADDR" envDefault:":8080"`
}

func DSN(c *Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName,
	)
}

func Load() (Config, error) {
	var cfg Config
	return cfg, env.Parse(&cfg)
}
