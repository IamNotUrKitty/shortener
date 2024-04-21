package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

func (s *LinksSuite) TestCreateLink() {
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
		// {
		// 	name: "Передан некорректный Content-type",
		// 	want: want{
		// 		code:        http.StatusBadRequest,
		// 		contentType: "text/plain; charset=UTF-8",
		// 		response:    "Неверный Content-type",
		// 	},
		// 	payload: payload{
		// 		body:        "https://ya.ru",
		// 		contentType: "application/json",
		// 	},
		// },
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

	for _, test := range tests {
		s.Run(test.name, func() {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.payload.body))
			request.Header.Add("Content-type", test.payload.contentType)

			w := httptest.NewRecorder()

			s.e.ServeHTTP(w, request)

			s.Equal(test.want.code, w.Code)

			s.Equal(test.want.response, w.Body.String())

			s.Equal(test.want.contentType, w.Header().Get("Content-type"))
		})
	}
}
