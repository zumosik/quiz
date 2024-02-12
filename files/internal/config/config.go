package config

import (
	firebase "firebase.google.com/go"
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"google.golang.org/api/option"
	"os"
	"time"
)

type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	GRPC           GRPCConfig `yaml:"grpc" env-required:"true"`
	StorageBucket  string     `yaml:"storage_bucket" env-required:"true"`
	DatabaseURL    string     `yaml:"database_url" env-required:"true"`
	StorageOptions option.ClientOption
	StorageCfg     *firebase.Config
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

func MustLoad() *Config {
	configPath, firebaseOptPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	cfg := MustLoadPath(configPath)

	opt := option.WithCredentialsFile(firebaseOptPath)
	config := &firebase.Config{
		DatabaseURL:   cfg.DatabaseURL,
		StorageBucket: cfg.StorageBucket,
	}

	cfg.StorageCfg = config
	cfg.StorageOptions = opt

	return cfg
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
func fetchConfigPath() (string, string) {
	var res, resFirebase string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.StringVar(&resFirebase, "firebase-config", "", "path to firebase json config file")
	flag.Parse()

	if res == "" {
		fmt.Printf("Getting config, path: %s \n", os.Getenv("CFG_PATH"))

		res = os.Getenv("CFG_PATH")
	}

	if resFirebase == "" {
		fmt.Printf("Getting Firebase cfg, path: %s \n", os.Getenv("FIREBASE_CFG_PATH"))

		resFirebase = os.Getenv("FIREBASE_CFG_PATH")
	}

	return res, resFirebase
}
