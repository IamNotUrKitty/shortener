package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
	"github.com/iamnoturkkitty/shortener/internal/echomiddleware"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateLink(c echo.Context) error {
	body, errBody := io.ReadAll(c.Request().Body)

	userID, err := echomiddleware.GetUser(c)
	if err != nil {
		return err
	}

	if errBody != nil {
		return c.String(http.StatusBadRequest, errBody.Error())
	}

	l, err := links.CreateLink(string(body), userID)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := h.repo.SaveLink(c.Request().Context(), *l); err != nil {
		if errors.Is(err, links.ErrLinkDuplicate) {
			return c.String(http.StatusConflict, h.baseAddress+"/"+l.Hash())
		}
		return c.String(http.StatusBadRequest, links.ErrLinkCreation.Error())
	}

	return c.String(http.StatusCreated, h.baseAddress+"/"+l.Hash())
}
