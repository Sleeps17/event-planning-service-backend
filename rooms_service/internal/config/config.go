package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env     string         `yaml:"env" env-default:"prod"`
	Server  GrpcConfig     `yaml:"server" env-required:"true"`
	Storage PostgresConfig `yaml:"storage"`
}

type GrpcConfig struct {
	Port    string        `yaml:"port" env-default:"4404"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

type PostgresConfig struct {
	Host     string        `yaml:"host" env-default:"localhost"`
	Port     string        `yaml:"port" env-default:"5432"`
	User     string        `yaml:"username" env-required:"true"`
	Password string        `yaml:"password" env-required:"true"`
	Database string        `yaml:"database" env-required:"true"`
	Timeout  time.Duration `yaml:"timeout"`
}

func (c *PostgresConfig) Url() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", c.User, c.Password, c.Host, c.Port, c.Database)
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		panic("config file does not exist:" + configPath)
	}

	return MustLoadByPath(configPath)
}

func MustLoadByPath(configPath string) *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); os.IsNotExist(err) {
		panic("can't read config: " + err.Error())
	}

	return &cfg
}
