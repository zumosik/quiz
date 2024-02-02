package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env                string        `yaml:"env" env-default:"local"`
	PostgresStorageURI string        `yaml:"postgres_storage_uri" env-required:"true"`
	GRPC               GRPCConfig    `yaml:"grpc"`
	TokenSecret        string        `yaml:"token_secret" env-required:"true"`
	TokenTTL           time.Duration `yaml:"token_ttl" env-default:"24h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		fmt.Printf("Getting config, path: %s \n", os.Getenv("CFG_PATH"))

		res = os.Getenv("CFG_PATH")
	}

	return res
}
