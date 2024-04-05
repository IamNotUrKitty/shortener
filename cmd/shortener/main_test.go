package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
				contentType: "",
			},
			payload: payload{
				body:        "https://ya.ru",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Передана некорректная ссылка для сокращения",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				response:    "Некорректный URL\n",
			},
			payload: payload{
				body:        "123",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Передан некорректный Content-type",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				response:    "Неверный Content-type\n",
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
				contentType: "text/plain; charset=utf-8",
				response:    "Некорректный URL\n",
			},
			payload: payload{
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.payload.body))
			request.Header.Add("Content-type", test.payload.contentType)

			w := httptest.NewRecorder()
			PostHandler(w, request)

			res := w.Result()

			defer res.Body.Close()

			resBody, _ := io.ReadAll(res.Body)

			assert.Equal(t, test.want.code, res.StatusCode)

			assert.Equal(t, test.want.response, string(resBody))

			assert.Equal(t, test.want.contentType, res.Header.Get("Content-type"))
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
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "Запросили ссылку которой не существует",
			payload: "123",
			want: want{
				code:        http.StatusBadRequest,
				response:    "URL не найден\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	urls["YBrECO"] = "https://ya.ru"

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+test.payload, nil)

			w := httptest.NewRecorder()
			GetHandler(w, request)

			res := w.Result()

			defer res.Body.Close()

			resBody, _ := io.ReadAll(res.Body)

			assert.Equal(t, test.want.code, res.StatusCode)

			if res.StatusCode == http.StatusBadRequest {
				assert.Equal(t, test.want.response, string(resBody))
			} else {
				assert.Equal(t, test.want.response, res.Header.Get("Location"))
			}

		})
	}
}

func TestRootHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/", nil)

	w := httptest.NewRecorder()
	RootHandler(w, request)

	res := w.Result()

	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	assert.Equal(t, "Метод не доступен\n", string(resBody))

	assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-type"))
}
