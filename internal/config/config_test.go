package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		want     string
		payload  string
		variable string
	}{
		{
			name:     "Прочитали переменную окружания",
			want:     "8888",
			payload:  "8888",
			variable: "SERVER_PORT",
		},
		{
			name:     "Получили значение по умолчанию если переменная не задана",
			want:     "8080",
			payload:  "",
			variable: "SERVER_PORT",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Unsetenv(test.variable)
			if test.payload != "" {
				os.Setenv(test.variable, test.payload)
			}
			env := getEnv(test.variable, "8080")
			assert.Equal(t, test.want, env)
		})
	}

}

func TestParseURL(t *testing.T) {
	tests := []struct {
		name       string
		payload    string
		defaultVal string
		error      bool
	}{
		{
			name:       "Корректно установили базовый URL",
			payload:    "http://ya.ru",
			defaultVal: "http://localhost:8080",
		},
		{
			name:       "Получили ошибку при не корректном URL",
			payload:    "123",
			defaultVal: "http://localhost:8080",
			error:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			URL := test.defaultVal

			res := parseURL(&URL, test.defaultVal)(test.payload)

			if test.error {
				assert.Error(t, res)

				return
			}

			assert.NoError(t, res)
			assert.Equal(t, test.payload, URL)
		})
	}

}

func TestParseAddress(t *testing.T) {
	tests := []struct {
		name       string
		payload    string
		defaultVal string
		error      bool
	}{
		{
			name:       "Корректный адрес запуска HTTP-сервера",
			payload:    "localhost:8081",
			defaultVal: "localhost:8080",
		},
		{
			name:       "Получили ошибку при не корректном адресе",
			payload:    "123",
			defaultVal: "localhost:8080",
			error:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			address := test.defaultVal

			res := parseAddress(&address, test.defaultVal)(test.payload)

			if test.error {
				assert.Error(t, res)

				return
			}

			assert.NoError(t, res)
			assert.Equal(t, test.payload, address)
		})
	}
}
