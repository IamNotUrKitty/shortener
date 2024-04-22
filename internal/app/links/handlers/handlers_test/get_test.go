package handlers_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
)

func (s *LinksSuite) TestGetLink() {
	link, _ := links.NewLink("https://ya.ru")

	err := s.repo.SaveLink(*link)

	s.Require().NoError(err)

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
				response:    links.ErrLinkNotFound.Error(),
				contentType: "text/plain; charset=UTF-8",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			request := httptest.NewRequest(http.MethodGet, "/"+test.payload, nil)

			w := httptest.NewRecorder()

			s.e.ServeHTTP(w, request)

			s.Equal(test.want.code, w.Code)

			if w.Code == http.StatusBadRequest {
				s.Equal(test.want.response, w.Body.String())
			} else {
				s.Equal(test.want.response, w.Header().Get("Location"))
			}
		})
	}

}
