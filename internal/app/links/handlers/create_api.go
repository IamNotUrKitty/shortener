package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
	"github.com/labstack/echo/v4"
)

type RequestDTO struct {
	URL string `json:"url"`
}

type ResponseDTO struct {
	Result string `json:"result"`
}

func (h *Handler) CreateLinkJSON(c echo.Context) error {
	var data RequestDTO

	body, errBody := io.ReadAll(c.Request().Body)

	if errBody != nil {
		return c.String(http.StatusBadRequest, errBody.Error())
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return c.String(http.StatusBadRequest, "Ошибка создания короткой ссылки")
	}

	l, err := links.NewLink(data.URL)

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
