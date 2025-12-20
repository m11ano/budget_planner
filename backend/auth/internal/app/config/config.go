package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BackendApp struct {
		Name    string `yaml:"name" env:"BACKEND_APP_NAME" env-default:"backend"`
		Version string `yaml:"version" env:"BACKEND_APP_VERSION" env-default:"1.0.0"`
		Base    struct {
			StartTimeoutSec int  `yaml:"start_timeout_sec" env:"BACKEND_APP_BASE_START_TIMEOUT_SEC" env-default:"10"`
			StopTimeoutSec  int  `yaml:"stop_timeout_sec" env:"BACKEND_APP_BASE_STOP_TIMEOUT_SEC" env-default:"2"`
			IsProd          bool `yaml:"is_prod" env:"BACKEND_APP_BASE_IS_PROD" env-default:"false"`
			UseFxLogger     bool `yaml:"use_fx_logger" env:"BACKEND_APP_BASE_USE_FX_LOGGER" env-default:"true"`
			UseLogger       bool `yaml:"use_logger" env:"BACKEND_APP_BASE_USE_LOGGER" env-default:"true"`
			LogSQLQueries   bool `yaml:"log_sql_queries" env:"BACKEND_APP_BASE_LOG_SQL_QUERIES" env-default:"true"`
			LogHTTP         bool `yaml:"log_http" env:"BACKEND_APP_BASE_LOG_HTTP" env-default:"true"`
		} `yaml:"base"`
		GRPC struct {
			Port       int  `yaml:"port" env:"BACKEND_APP_GRPC_PORT" env-default:"50051"`
			LogQueries bool `yaml:"log_queries" env:"BACKEND_APP_GRPC_LOG_QUERIES" env-default:"false"`
		} `yaml:"grpc"`
	} `yaml:"backend_app"`
	Postgres struct {
		MaxAttempts         int    `yaml:"max_attempts" env:"POSTGRES_MAX_ATTEMPTS" env-default:"3"`
		AttemptSleepSeconds int    `yaml:"attempt_sleep_seconds" env:"POSTGRES_ATTEMPT_SLEEP_SECONDS" env-default:"1"`
		MigrationsPath      string `yaml:"migrations_path" env:"POSTGRES_MIGRATIONS_PATH" env-default:"migrations"`
		Master              struct {
			DSN string `yaml:"dsn" env:"POSTGRES_MASTER_DSN"`
		} `yaml:"master"`
	} `yaml:"postgres"`
	Auth struct {
		AccessTokenLifetimeSec int `yaml:"access_token_lifetime_sec" env:"AUTH_ACCESS_TOKEN_LIFETIME_SEC" env-default:"300"`
		// nolint
		RefreshTokenLifetimeHrs int    `yaml:"refresh_token_lifetime_hours" env:"AUTH_REFRESH_TOKEN_LIFETIME_HRS" env-default:"720"`
		JwtAccessSecret         string `yaml:"jwt_access_secret" env:"AUTH_JWT_ACCESS_SECRET"`
		JwtRefreshSecret        string `yaml:"jwt_refresh_secret" env:"AUTH_JWT_REFRESH_SECRET"`
	}
}

func LoadConfig(files ...string) Config {
	var Config Config

	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			err := cleanenv.ReadConfig(file, &Config)
			if err != nil {
				log.Fatal("config file error ", err)
			}
		}
	}

	return Config
}
