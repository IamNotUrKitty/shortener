package handlers

import (
	"net/http"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetLink(c echo.Context) error {
	hash := c.Param("hash")

	link, err := h.repo.GetLink(hash)
	if err != nil {
		return c.String(http.StatusBadRequest, links.ErrLinkNotFound.Error())
	} else {
		return c.Redirect(http.StatusTemporaryRedirect, link.URL())
	}
}
