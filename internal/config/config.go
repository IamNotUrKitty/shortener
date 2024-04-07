package config

import "flag"

var (
	Address     *string
	BaseAddress *string
)

func init() {
	Address = flag.String("a", "localhost:8080", "Адрес запуска HTTP-сервера")
	BaseAddress = flag.String("b", "localhost:8080", "Базовый адрес результирующего сокращённого URL")
}
