package config

import "flag"

type Config struct {
	DSN string
}

func GetConfig() *Config {
	c := &Config{}
	flag.StringVar(&c.DSN, "dsn", "", "connection string")

	flag.Parse()

	return c
}
