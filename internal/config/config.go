package config

import (
	"errors"
	"flag"
	"net/url"
	"os"
	"strings"
)

var (
	// Адрес запуска HTTP-сервера
	Address string
	// Базовый адрес результирующего сокращённого URL
	BaseAddress string
)

func parseAddress(configField *string, defaultValue string) func(string) error {
	return func(address string) error {
		hp := strings.Split(address, ":")
		if len(hp) != 2 {
			return errors.New("need address in a form host:port")
		}

		if *configField == defaultValue {
			*configField = address
		}

		return nil
	}
}

func parseURL(configField *string, defaultValue string) func(string) error {
	return func(address string) error {
		_, errURL := url.ParseRequestURI(address)
		if errURL != nil {
			return errors.New("need base address in a valid URL form")
		}

		if *configField == defaultValue {
			*configField = address
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

func init() {
	defaultAddress := "localhost:8080"
	defaultBaseAddress := "http://localhost:8080"

	flag.Func("a", "Адрес запуска HTTP-сервера", parseAddress(&Address, defaultAddress))
	flag.Func("b", "Базовый адрес результирующего сокращённого URL", parseURL(&BaseAddress, defaultBaseAddress))

	Address = getEnv("SERVER_ADDRESS", defaultAddress)
	BaseAddress = getEnv("BASE_URL", defaultBaseAddress)

	flag.Parse()
}
