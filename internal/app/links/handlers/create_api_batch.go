package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
	"github.com/labstack/echo/v4"
)

type RequestBatchDTO struct {
	URL           string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
}

type ResponseBatchDTO struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (h *Handler) CreateLinkBatch(c echo.Context) error {
	var data []RequestBatchDTO
	var response []ResponseBatchDTO
	var linkArr []links.Link

	body, errBody := io.ReadAll(c.Request().Body)

	if errBody != nil {
		return c.String(http.StatusBadRequest, errBody.Error())
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return c.String(http.StatusBadRequest, links.ErrLinkCreation.Error())
	}

	for _, i := range data {
		l, err := links.CreateLink(i.URL)
		if err != nil {
			return c.String(http.StatusBadRequest, links.ErrLinkCreation.Error())
		}

		linkArr = append(linkArr, *l)

		response = append(response, ResponseBatchDTO{CorrelationID: i.CorrelationID, ShortURL: h.baseAddress + "/" + l.Hash()})
	}

	if err := h.repo.SaveLinkBatch(linkArr); err != nil {
		return c.String(http.StatusBadRequest, links.ErrLinkCreation.Error())
	}

	return c.JSON(http.StatusCreated, response)
}
