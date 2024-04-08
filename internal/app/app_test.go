package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestPostHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	type payload struct {
		body        string
		contentType string
	}

	tests := []struct {
		name    string
		want    want
		payload payload
	}{
		{
			name: "Успешно cоздали короткую ссылку",
			want: want{
				code:        http.StatusCreated,
				response:    "http://localhost:8080/YBrECO",
				contentType: "text/plain; charset=UTF-8",
			},
			payload: payload{
				body:        "https://ya.ru",
				contentType: "text/plain; charset=UTF-8",
			},
		},
		{
			name: "Передана некорректная ссылка для сокращения",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				response:    "Некорректный URL",
			},
			payload: payload{
				body:        "123",
				contentType: "text/plain; charset=UTF-8",
			},
		},
		{
			name: "Передан некорректный Content-type",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				response:    "Неверный Content-type",
			},
			payload: payload{
				body:        "https://ya.ru",
				contentType: "application/json",
			},
		},
		{
			name: "Передано пустое тело запроса",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				response:    "Некорректный URL",
			},
			payload: payload{
				body:        "",
				contentType: "text/plain; charset=UTF-8",
			},
		},
	}

	e := echo.New()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.payload.body))
			request.Header.Add("Content-type", test.payload.contentType)

			w := httptest.NewRecorder()

			c := e.NewContext(request, w)

			makePostHandler("http://localhost:8080")(c)

			assert.Equal(t, test.want.code, w.Code)

			assert.Equal(t, test.want.response, w.Body.String())

			assert.Equal(t, test.want.contentType, w.Header().Get("Content-type"))
		})
	}
}

func TestGetHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		want    want
		payload string
	}{
		{
			name:    "Успешно получили ссылку",
			payload: "YBrECO",
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    "https://ya.ru",
				contentType: "text/plain; charset=UTF-8",
			},
		},
		{
			name:    "Запросили ссылку которой не существует",
			payload: "123",
			want: want{
				code:        http.StatusBadRequest,
				response:    "URL не найден",
				contentType: "text/plain; charset=UTF-8",
			},
		},
	}

	urls["YBrECO"] = "https://ya.ru"

	e := echo.New()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)

			w := httptest.NewRecorder()

			c := e.NewContext(request, w)

			if test.payload != "" {
				c.SetPath("/:hash")
				c.SetParamNames("hash")
				c.SetParamValues(test.payload)
			}

			getHandler(c)

			assert.Equal(t, test.want.code, w.Code)

			if w.Code == http.StatusBadRequest {
				assert.Equal(t, test.want.response, w.Body.String())
			} else {
				assert.Equal(t, test.want.response, w.Header().Get("Location"))
			}
		})
	}
}
