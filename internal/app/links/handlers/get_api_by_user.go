package handlers

import (
	"net/http"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
	"github.com/iamnoturkkitty/shortener/internal/echomiddleware"
	"github.com/labstack/echo/v4"
)

type ResponseByUserDTO struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

func (h *Handler) GetLinksByUserID(c echo.Context) error {
	var res = []ResponseByUserDTO{}

	_, err := c.Cookie(echomiddleware.CookieName)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	userID, err := echomiddleware.GetUser(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	linkList, err := h.repo.GetLinkByUserID(c.Request().Context(), userID)

	if err != nil {
		return c.String(http.StatusBadRequest, links.ErrLinkNotFound.Error())
	} else {
		if len(linkList) == 0 {
			return c.JSON(http.StatusNoContent, res)
		}

		for _, l := range linkList {
			res = append(res, ResponseByUserDTO{OriginalURL: l.URL(), ShortURL: h.baseAddress + "/" + l.Hash()})
		}

		return c.JSON(http.StatusOK, res)
	}
}
