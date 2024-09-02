package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type RabbitMq struct {
	Addr string `yaml:"addr" env-required:"true"`
}

type GRPC struct {
	Port    int           `yaml:"port" env-default:"4041"`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
}

type Config struct {
	Env      string   `yaml:"env" env-default:"local"`
	GRpc     GRPC     `yaml:"grpc_app" env-required:"true"`
	Storage  Storage  `yaml:"storage" env-required:"true"`
	RabbitMq RabbitMq `yaml:"rabbitmq" env-required:"true"`
}

type Storage struct {
	ConnectionString string `yaml:"connection_string" env-required:"true"`
	MigrationPath    string `yaml:"migration_path" env-required:"true"`
}

func MustLoad(configPath string) *Config {
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
