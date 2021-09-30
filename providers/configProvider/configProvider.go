package configProvider

import (
	"github.com/bms/providers"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func NewConfigProvider() providers.ConfigProvider {
	return &Config{}
}

func (c *Config) Read() error {
	err := envconfig.Process("", c)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	return nil
}

func (c *Config) GetServerPort() string {
	if c == nil {
		return ""
	}
	return c.Port
}

func (c *Config) GetJWTKey() string {
	if c == nil {
		return ""
	}
	return c.JWTKey
}

func (c *Config) GetString(key string) string {
	return os.Getenv(key)
}

func (c *Config) GetInt(key string) int {
	intVal, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		logrus.Errorf("Error getting int config %v", err)
	}
	return intVal
}

func (c *Config) GetAny(key string) interface{} {
	return os.Getenv(key)
}

func (c *Config) GetPSQLConnectionString() string {
	return c.DBString
}

func (c *Config) GetPSQLMaxConnection() int {
	return c.MaxConnections
}

func (c *Config) GetPSQLMaxIdleConnection() int {
	return c.MaxIdleConnections
}
