package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config parse
type Config struct {
	Env     string `yaml:"env" env-default:"local" env-required:"true"`
	Storage `yaml:"storage"`
	Network `yaml:"network"`
}

// Network parse
type Network struct {
	Address string `yaml:"address"  env-default:"localhost"`
	Port    int    `yaml:"port"  env-default:"5551"`
}

// Storage parse
type Storage struct {
	StoragePath string `yaml:"storage_path" env-required:"true"`
	Name        string `yaml:"db_name"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
}

// MustLoad parse config and returns Config struct object
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	// проверка на существование файла
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Configuration file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cant read from config, error: %s", err)
	}

	return &cfg
}
