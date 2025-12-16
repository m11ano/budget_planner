package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Global struct {
		ServerIP string `yaml:"server_ip" env:"GLOBAL_SERVER_IP"`
	} `yaml:"global"`
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
		HTTP struct {
			Port             int      `yaml:"port" env:"BACKEND_APP_HTTP_PORT" env-default:"8080"`
			Prefix           string   `yaml:"prefix" env:"BACKEND_APP_HTTP_PREFIX" env-default:""`
			UnderProxy       bool     `yaml:"under_proxy" env:"BACKEND_APP_HTTP_UNDER_PROXY" env-default:"false"`
			StopTimeoutSec   int      `yaml:"stop_timeout_sec" env:"BACKEND_APP_HTTP_STOP_TIMEOUT_SEC" env-default:"3"`
			CorsAllowOrigins []string `yaml:"cors_allow_origins" env:"BACKEND_APP_HTTP_CORS_ALLOW_ORIGINS" env-default:""`
		} `yaml:"http"`
	} `yaml:"backend_app"`
	Auth struct {
		JwtAccessSecret string `yaml:"jwt_access_secret" env:"AUTH_JWT_ACCESS_SECRET"`
	}
	GRPC struct {
		Auth struct {
			Addr         string `yaml:"addr" env:"GRPC_AUTH_ADDR"`
			RetriesCount int    `yaml:"retries_count" env:"GRPC_AUTH_RETRIES_COUNT" env-default:"5"`
			Timeout      string `yaml:"timeout" env:"GRPC_AUTH_TIMEOUT" env-default:"60s"`
		}
	}
}

func LoadConfig(files ...string) Config {
	var Config Config

	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			err := cleanenv.ReadConfig(file, &Config)
			if err != nil {
				log.Fatal("config file error", err)
			}
		}
	}

	return Config
}
