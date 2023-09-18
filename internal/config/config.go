package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

const (
	Path = "CONFIG_PATH"
)

type Config struct {
	StorageConnectString string `yaml:"storage_connect_string"`
	HTTPServer           `yaml:"http_server"`
	Broker               `yaml:"broker"`
}

type HTTPServer struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port string `yaml:"port" env-default:"8080"`
}

type Broker struct {
	ClusterID string `yaml:"cluster_id" env-default:"test-cluster"`
	ClientID  string `yaml:"client_id" env-default:"test-client"`
	URL       string `yaml:"url" env-default:"nats://localhost:4222"`
}

func MustLoadConfig() (*Config, error) {
	const fn = "config.MustLoadConfig"
	err := godotenv.Load()
	if err != nil {
		slog.Error(fn, slog.String("failed to load env file", err.Error()))
		return nil, err
	}
	slog.Info("env file is loaded")

	configPath := os.Getenv(Path)
	if configPath == "" {
		slog.Error(fn, slog.String("config path is not set", err.Error()))
		return nil, err
	}
	slog.Info("config path is set")

	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		slog.Error(fn, slog.String("config file does not exist", err.Error()))
		return nil, err
	}
	slog.Info("config file exists")

	var httpServer HTTPServer
	var broker Broker
	var cnf = &Config{
		HTTPServer: httpServer,
		Broker:     broker,
	}

	if err = cleanenv.ReadConfig(configPath, cnf); err != nil {
		slog.Error(fn, slog.String("failed to read config", err.Error()))
		return nil, err
	}
	slog.Info("config is loaded successfully")

	return cnf, nil
}
