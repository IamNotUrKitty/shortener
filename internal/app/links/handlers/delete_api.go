package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
	"github.com/iamnoturkkitty/shortener/internal/echomiddleware"
	"github.com/labstack/echo/v4"
)

type RequestDeleteDTO []string

func (h *Handler) DeleteLinkByUserID(c echo.Context) error {
	var data = RequestDeleteDTO{}

	_, err := c.Cookie(echomiddleware.CookieName)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	userID, err := echomiddleware.GetUser(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	body, errBody := io.ReadAll(c.Request().Body)
	if errBody != nil {
		return c.String(http.StatusBadRequest, errBody.Error())
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return c.String(http.StatusBadRequest, links.ErrLinkCreation.Error())
	}

	h.repo.DeleteLinkBatch(c.Request().Context(), data, userID)

	return c.JSON(http.StatusAccepted, data)

}
