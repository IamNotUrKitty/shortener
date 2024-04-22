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
		return c.String(http.StatusBadRequest, links.ErrLinkCreation.Error())
	}

	l, err := links.NewLink(data.URL)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := h.repo.SaveLink(*l); err != nil {
		return c.String(http.StatusBadRequest, links.ErrLinkCreation.Error())
	}

	return c.JSON(http.StatusCreated, ResponseDTO{Result: h.baseAddress + "/" + l.Hash()})
}
