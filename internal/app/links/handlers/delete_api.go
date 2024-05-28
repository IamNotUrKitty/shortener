package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
	"github.com/iamnoturkkitty/shortener/internal/echomiddleware"
	"github.com/labstack/echo/v4"
)

func deleteTaskGenerator(hashes []string, userID int, queue chan<- links.DeleteLinkTask) error {
	go func() {
		for _, val := range hashes {
			queue <- links.DeleteLinkTask{Hash: val, UserID: userID}
		}
	}()

	return nil
}

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

	fmt.Println(userID)

	body, errBody := io.ReadAll(c.Request().Body)
	if errBody != nil {
		return c.String(http.StatusBadRequest, errBody.Error())
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return c.String(http.StatusBadRequest, links.ErrLinkCreation.Error())
	}

	deleteTaskGenerator(data, userID, h.deleteQueue)

	return c.JSON(http.StatusAccepted, data)

}
