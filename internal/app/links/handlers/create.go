package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
	"github.com/labstack/echo"
)

func (h *Handler) CreateLink(c echo.Context) error {
	// Валидация на сontent-type
	if strings.ToLower(c.Request().Header.Get("Content-type")) != "text/plain; charset=utf-8" {
		return c.String(http.StatusBadRequest, "Неверный Content-type")
	}

	body, errBody := io.ReadAll(c.Request().Body)

	if errBody != nil {
		return c.String(http.StatusBadRequest, errBody.Error())
	}

	l, err := links.NewLink(string(body))

	if err != nil {
		return c.String(http.StatusBadRequest, "Ошибка создания короткой ссылки")
	}

	h.repo.SaveLink(*l)

	return c.String(http.StatusCreated, "http:localhost:8080/"+l.Hash())
}