package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func (h *Handler) GetLink(c echo.Context) error {
	hash := c.Param("hash")

	link, err := h.repo.GetLink(hash)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, link.Url())
	} else {
		return c.String(http.StatusBadRequest, "URL не найден")
	}
}
