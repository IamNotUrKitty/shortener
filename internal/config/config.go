package config

import (
	"errors"
	"flag"
	"net/url"
	"strings"
)

var (
	// Адрес запуска HTTP-сервера
	Address string = "localhost:8080"
	// Базовый адрес результирующего сокращённого URL
	BaseAddress string = "http://localhost:8080"
)

func parseAddress(addressP *string) func(string) error {
	return func(address string) error {
		hp := strings.Split(address, ":")
		if len(hp) != 2 {
			return errors.New("need address in a form host:port")
		}

		*addressP = address

		return nil
	}
}

func parseURL(addressP *string) func(string) error {
	return func(address string) error {
		_, errURL := url.ParseRequestURI(address)
		if errURL != nil {
			return errors.New("need base address in a valid URL form")
		}

		*addressP = address

		return nil
	}
}

func init() {
	flag.Func("a", "Адрес запуска HTTP-сервера", parseAddress(&Address))

	flag.Func("b", "Базовый адрес результирующего сокращённого URL", parseURL(&BaseAddress))
}
