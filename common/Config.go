package common

import "os"

type Config struct {
	apiBaseUrl   string
	dbConnString string
}

func (c *Config) GetApiBaseUrl() string {
	return c.apiBaseUrl
}

func (c *Config) GetDbConnString() string {
	return c.dbConnString
}

func NewEnvConfig() *Config {
	c := Config{
		apiBaseUrl:   os.Getenv("LF_API_BASE_URL"),
		dbConnString: os.Getenv("LF_DB_CONN_STRING"),
	}

	return &c
}
