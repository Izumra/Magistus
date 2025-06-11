package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Bot Bot `yaml:"bot"`
	Db  Db  `yaml:"db"`
}

type Bot struct {
	Token string `yaml:"token"`
}

type Db struct {
	Path string `yaml:"path"`
}

func MustLoad() *Config {
	path := fetchConfigByPath()
	if path == "" {
		panic("Failed to loading config: Config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {
	stream, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = yaml.Unmarshal(stream, &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

func fetchConfigByPath() string {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to the config file")
	flag.Parse()

	if configPath == "" {
		godotenv.Load(".env")
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}
