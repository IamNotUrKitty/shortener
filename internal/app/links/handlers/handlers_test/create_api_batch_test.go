package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/iamnoturkkitty/shortener/internal/app/links/handlers"
	"github.com/labstack/echo/v4"
)

func (s *LinksSuite) TestCreateApiLinkBatch() {
	payload := []handlers.RequestBatchDTO{
		{URL: "https://ya.ru", CorrelationID: "1"},
		{URL: "https://test.ru", CorrelationID: "2"},
		{URL: "https://testfoo.ru", CorrelationID: "3"},
	}

	reqBody, _ := json.Marshal(payload)

	request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(reqBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	w := httptest.NewRecorder()

	s.e.ServeHTTP(w, request)

	s.Require().Equal(http.StatusCreated, w.Code)

	var resp []handlers.ResponseBatchDTO

	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)

	s.Require().Equal(len(payload), len(resp))
}
