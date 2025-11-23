package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string `yaml:"env" env-default:"local"`
	HTTPServer    `yaml:"http_server"`
	DB            `yaml:"db"`
	LogLevel      string `yaml:"log_level" env-default:"info"`
	EnableSwagger bool   `yaml:"enable_swagger" env-default:"true"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:":8080"`
}

type DB struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     uint   `yaml:"port" env-default:"5432"`
	Username string `yaml:"username" env-default:"postgres"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"db_name" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
