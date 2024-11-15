package config

import (
	"github.com/gookit/slog"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	GrpcPort int      `yaml:"grpc_port"`
	Database postgres `yaml:"postgres"`
	Redis    redis    `yaml:"redis"`
	Token    token    `yaml:"token"`
	Gateway  gateway  `yaml:"gateway"`
}

type postgres struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"db_name"`
}

type redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	RedisDB  int    `yaml:"db"`
}

type token struct {
	TKey  string `yaml:"T_KEY"`
	RTKey string `yaml:"RT_KEY"`
}

type gateway struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func LoadConfig() *Config {
	var cfg *Config
	wd, err := os.Getwd()
	if err != nil {
		slog.Fatalf("Failed to get working directory: %v", err)
	}
	configPath := filepath.Join(wd, "config.yml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		slog.Fatalf("Failed to read configuration file: %v", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		slog.Fatalf("Failed to unmarshal YAML data to config: %v", err)
	}
	return cfg
}
