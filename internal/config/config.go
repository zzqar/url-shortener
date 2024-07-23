package config

import (
	//go get -u github.com/ilyakaznacheev/cleanenv - установка пакета для чтения конфигурационных файлов
	"github.com/ilyakaznacheev/cleanenv"

	"log"
	"os"
	"time"
)

type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-default:"prod"`
	StoragePath string     `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
}

/*
yaml, env - парсинг из файлов с таким форматом
env-default - устанавливает значение переменной окружения, если оно не указано в конфигурационном файле
env-required - не соберет приложение, если в конфигурационном файле не указано значение
*/

type HTTPServer struct {
	Address      string        `yaml:"address" env:"ADDRESS" env-default:"localhost:8080"`
	Timeout      time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"5s"`
	IddleTimeout time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"60s"`
}

func MustLoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	//check if exist file
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file %s not found", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	return &cfg
}
