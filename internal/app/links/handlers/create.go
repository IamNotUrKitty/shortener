package handlers

import (
	"compress/gzip"
	"io"
	"net/http"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
	"github.com/labstack/echo/v4"
)

func GetBody(c echo.Context) ([]byte, error) {
	var reader io.Reader

	if c.Request().Header.Get(echo.HeaderContentEncoding) == `gzip` {
		gz, err := gzip.NewReader(c.Request().Body)
		if err != nil {
			return nil, err
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = c.Request().Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (h *Handler) CreateLink(c echo.Context) error {
	// Валидация на сontent-type
	// if strings.ToLower(c.Request().Header.Get(echo.HeaderContentType)) != "text/plain; charset=utf-8" {
	// 	return c.String(http.StatusBadRequest, "Неверный Content-type")
	// }

	body, errBody := GetBody(c)

	if errBody != nil {
		return c.String(http.StatusBadRequest, errBody.Error())
	}

	l, err := links.NewLink(string(body))

	if err != nil {
		// TODO: Убрать сравнение со стрингой
		if err.Error() == "bad url" {
			return c.String(http.StatusBadRequest, "Некорректный URL")
		}
		return c.String(http.StatusBadRequest, "Ошибка создания короткой ссылки")
	}

	h.repo.SaveLink(*l)

	return c.String(http.StatusCreated, h.baseAddress+"/"+l.Hash())
}
