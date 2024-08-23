package config

import (
	"flag"
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

func MustLoad() *Config {
	var cfg Config

	configPath := fetchConfig()
	if configPath == "" {
		panic("the path to the config is not specified")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file is not exists")
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("problem reading config")
	}

	return &cfg
}

func fetchConfig() string {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to config")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}
