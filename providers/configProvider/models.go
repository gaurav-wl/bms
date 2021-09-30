package configProvider

type Config struct {
	PSQLConfig
	HTTPKubeConfig
	AUTHKubeConfig
}

type PSQLConfig struct {
	DBString           string `envconfig:"PSQL_DB_URL"`
	MaxConnections     int    `envconfig:"PSQL_DB_MAX_CONNECTIONS"`
	MaxIdleConnections int    `envconfig:"PSQL_DB_MAX_IDLE_CONNECTIONS"`
}

type HTTPKubeConfig struct {
	Port string
}

type AUTHKubeConfig struct {
	JWTKey string `envconfig:"JWT_KEY"`
}
