package config

import (
	"errors"
	"flag"
	"net/url"
	"os"
	"strings"
)

func parseAddress(configField *string, defaultValue string) func(string) error {
	return func(value string) error {
		hp := strings.Split(value, ":")
		if len(hp) != 2 {
			return errors.New("need address in a form host:port")
		}

		if *configField == defaultValue {
			*configField = value
		}

		return nil
	}
}

func parseURL(configField *string, defaultValue string) func(string) error {
	return func(value string) error {
		_, errURL := url.ParseRequestURI(value)
		if errURL != nil {
			return errors.New("need base address in a valid URL form")
		}

		if *configField == defaultValue {
			*configField = value
		}

		return nil
	}
}

func parseStorageFile(configField *string, defaultValue string) func(string) error {
	return func(value string) error {
		if *configField == defaultValue {
			*configField = value
		}

		return nil
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.ToLower(value)
	}

	return fallback
}

type Config struct {
	// Адрес запуска HTTP-сервера
	Address string
	// Базовый адрес результирующего сокращённого URL
	BaseAddress string
	StorageFile string
}

func GetConfig() *Config {
	var cfg Config
	defaultAddress := "localhost:8080"
	defaultBaseAddress := "http://localhost:8080"
	defaultStorageFile := ""

	flag.Func("a", "Адрес запуска HTTP-сервера", parseAddress(&cfg.Address, defaultAddress))
	flag.Func("b", "Базовый адрес результирующего сокращённого URL", parseURL(&cfg.BaseAddress, defaultBaseAddress))
	flag.Func("f", "Базовый адрес результирующего сокращённого URL", parseStorageFile(&cfg.StorageFile, defaultStorageFile))

	cfg.Address = getEnv("SERVER_ADDRESS", defaultAddress)
	cfg.BaseAddress = getEnv("BASE_URL", defaultBaseAddress)
	cfg.StorageFile = getEnv("FILE_STORAGE_PATH", defaultStorageFile)

	flag.Parse()

	return &cfg
}
