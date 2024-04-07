package config

import (
	"errors"
	"flag"
	"strings"
)

var (
	Address     string = "localhost:8080"
	BaseAddress string = "localhost:8080"
)

func parseAddress(addressP *string) func(string) error {
	return func(address string) error {
		hp := strings.Split(address, ":")
		if len(hp) != 2 {
			return errors.New("Need address in a form host:port")
		}

		*addressP = address

		return nil
	}
}

func init() {
	flag.Func("a", "Адрес запуска HTTP-сервера", parseAddress(&Address))

	flag.Func("b", "Базовый адрес результирующего сокращённого URL", parseAddress(&BaseAddress))
}
