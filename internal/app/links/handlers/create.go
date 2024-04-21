package handlers

import (
	"io"
	"net/http"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateLink(c echo.Context) error {
	body, errBody := io.ReadAll(c.Request().Body)

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
