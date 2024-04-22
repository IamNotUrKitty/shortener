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
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := h.repo.SaveLink(*l); err != nil {
		return c.String(http.StatusBadRequest, links.ErrLinkCreation.Error())
	}

	return c.String(http.StatusCreated, h.baseAddress+"/"+l.Hash())
}
