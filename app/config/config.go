package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	RedisHostPort string `env:"REDIS_HOSTPORT" env-default:"localhost:6379"`
	ServiceHost   string `env:"SERVICE_HOST" env-default:"localhost:8083"`
	StreamName    string `env:"REDIS_STREAM_NAME" env-default:"my_stream"`
}

func GetConfig() Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal("Can not read config", err)
	}
	return cfg
}
