package main

import (
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/sqids/sqids-go"
)

// мапа для хранения значений [хэш]урл
var urls = make(map[string]string)

// библиотека для генерации хэшей
var s, _ = sqids.New()

// Генерация хэша из массива байт
func makeHash(byteURL []byte) (string, error) {
	d := make([]uint64, len(byteURL))

	for i, b := range byteURL {
		d[i] = uint64(b)
	}

	hash, err := s.Encode(d)

	return hash[:6], err
}

// Обработчик POST запросов
func PostHandler(c echo.Context) error {
	// Валидация на сontent-type
	if c.Request().Header.Get("Content-type") != "text/plain; charset=UTF-8" {
		return c.String(http.StatusBadRequest, "Неверный Content-type")
	}

	body, errBody := io.ReadAll(c.Request().Body)

	if errBody != nil {
		return c.String(http.StatusBadRequest, errBody.Error())
	}

	// Валидация корректности URL
	_, errURL := url.ParseRequestURI(string(body))

	if errURL != nil {
		return c.String(http.StatusBadRequest, "Некорректный URL")
	}

	hash, errHash := makeHash(body)
	if errHash != nil {
		return c.String(http.StatusBadRequest, "Ошибка создания короткой ссылки")
	}

	urls[hash] = string(body)

	return c.String(http.StatusCreated, "http://localhost:8080/"+hash)
}

// Обработчик GET запросов
func GetHandler(c echo.Context) error {
	hash := c.Param("hash")

	url, ok := urls[hash]
	if ok {
		return c.Redirect(http.StatusTemporaryRedirect, url)
	} else {
		return c.String(http.StatusBadRequest, "URL не найден")
	}
}

func main() {
	e := echo.New()

	e.GET("/:hash", GetHandler)
	e.POST("/", PostHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
