package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	Addr     string        `yaml:"addr" env-required:"true"`
	Services ServicesAddrs `yaml:"services" env-required:"true"`
	Timeout  time.Duration `yaml:"timeout" env-default:"15s"`
}

type ServicesAddrs struct {
	Files string `yaml:"files" env-required:"true"`
	Users string `yaml:"users" env-required:"true"`
}

func MustGetCfgPath() string {
	s := os.Getenv("CONFIG_PATH")
	if s == "" {
		flag.StringVar(&s, "config-path", "", "path to .yml config")
	}
	if s == "" {
		panic("empty cfg path")
	}

	return s
}

func MustGetConfig() Config {
	s := MustGetCfgPath()

	f, err := os.Open(s)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg Config

	if err := cleanenv.ParseYAML(f, &cfg); err != nil {
		panic(err)
	}

	return cfg
}
