package handlers

import (
	"net/http"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
	"github.com/labstack/echo/v4"
)

type RequestDTO struct {
	Url string `json:"url"`
}

type ResponseDTO struct {
	Result string `json:"result"`
}

func (h *Handler) CreateLinkJSON(c echo.Context) error {
	var data RequestDTO

	if err := c.Bind(&data); err != nil {
		return err
	}

	l, err := links.NewLink(data.Url)

	if err != nil {
		// TODO: Убрать сравнение со стрингой
		if err.Error() == "bad url" {
			return c.String(http.StatusBadRequest, "Некорректный URL")
		}
		return c.String(http.StatusBadRequest, "Ошибка создания короткой ссылки")
	}

	if err := h.repo.SaveLink(*l); err != nil {
		return c.String(http.StatusBadRequest, "Ошибка создания короткой ссылки")
	}

	return c.JSON(http.StatusCreated, ResponseDTO{Result: h.baseAddress + "/" + l.Hash()})
}
