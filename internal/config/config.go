package config

import (
	"os"
)

type Config struct {
	Host string
	DSN  string
}

func GetConfig() *Config {

	c := &Config{}
	c.Host, _ = os.LookupEnv("EMPLOYEES_HOST")
	c.DSN, _ = os.LookupEnv("EMPLOYEES_DSN")

	return c
}
