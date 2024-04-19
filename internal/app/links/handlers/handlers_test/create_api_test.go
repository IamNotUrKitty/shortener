package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/iamnoturkkitty/shortener/internal/app/links/handlers"
	"github.com/labstack/echo/v4"
)

func (s *LinksSuite) TestCreateApiLink() {
	payload := handlers.RequestDTO{URL: "https://ya.ru"}

	reqBody, _ := json.Marshal(payload)

	request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(reqBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	w := httptest.NewRecorder()

	s.e.ServeHTTP(w, request)

	s.Require().Equal(http.StatusCreated, w.Code)

	var resp handlers.ResponseDTO

	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)
}
