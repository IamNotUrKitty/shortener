package main

import (
	"encoding/binary"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sqids/sqids-go"
)

func makeMainHandler() http.HandlerFunc {
	urls := make(map[string]string)
	s, _ := sqids.New()

	return func(res http.ResponseWriter, req *http.Request) {

		if req.Method == http.MethodPost {
			if req.Header.Get("Content-type") != "text/plain; charset=utf-8" {
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			body, errBody := io.ReadAll(req.Body)

			if errBody != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			_, errURL := url.ParseRequestURI(string(body))

			if errURL != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			defer req.Body.Close()

			// len < 7 problems
			data := binary.BigEndian.Uint64(body)

			hash, _ := s.Encode([]uint64{data})

			urls[hash] = string(body)

			res.WriteHeader(http.StatusCreated)

			res.Write([]byte("http://localhost:8080/" + hash))

			return
		}

		if req.Method == http.MethodGet {
			hash := strings.TrimPrefix(req.URL.Path, "/")
			url, ok := urls[hash]
			if ok {
				http.Redirect(res, req, url, http.StatusTemporaryRedirect)
			} else {
				res.WriteHeader(http.StatusBadRequest)
			}

			return
		}

		res.WriteHeader(http.StatusBadRequest)

	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", makeMainHandler())

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err)
	}
}
