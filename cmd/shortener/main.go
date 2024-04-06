package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sqids/sqids-go"
)

// мапа для хранения значений [хэш]урл
var urls = make(map[string]string)

// библиотека для генерации хэшей
var s, _ = sqids.New()

// Генерация хэша из массива байт
func makeHash(byteURL []byte) (string, error) {
	d := make([]uint64, len(byteURL))

	for i, b := range byteURL {
		d[i] = uint64(b)
	}

	hash, err := s.Encode(d)

	return hash[:6], err
}

// Обработчик POST запросов
func PostHandler(res http.ResponseWriter, req *http.Request) {
	// Валидация на сontent-type
	if req.Header.Get("Content-type") != "text/plain; charset=utf-8" {
		http.Error(res, "Неверный Content-type", http.StatusBadRequest)
		return
	}

	body, errBody := io.ReadAll(req.Body)

	if errBody != nil {
		http.Error(res, errBody.Error(), http.StatusBadRequest)
		return
	}

	// Валидация корректности URL
	_, errURL := url.ParseRequestURI(string(body))

	if errURL != nil {
		http.Error(res, "Некорректный URL", http.StatusBadRequest)
		return
	}

	hash, errHash := makeHash(body)
	if errHash != nil {
		http.Error(res, "Ошибка создания короткой ссылки", http.StatusBadRequest)
	}

	urls[hash] = string(body)

	res.WriteHeader(http.StatusCreated)

	res.Write([]byte("http://localhost:8080/" + hash))
}

// Обработчик GET запросов
func GetHandler(res http.ResponseWriter, req *http.Request) {
	hash := strings.TrimPrefix(req.URL.Path, "/")
	url, ok := urls[hash]
	if ok {
		http.Redirect(res, req, url, http.StatusTemporaryRedirect)
	} else {
		http.Error(res, "URL не найден", http.StatusBadRequest)
	}
}

// Обработчик входящего запроса
func RootHandler(res http.ResponseWriter, req *http.Request) {
	switch method := req.Method; method {
	case http.MethodPost:
		PostHandler(res, req)
	case http.MethodGet:
		GetHandler(res, req)
	default:
		// В случае метода который не обрабатываем возвращаем ошибку
		http.Error(res, "Метод не доступен", http.StatusBadRequest)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", RootHandler)

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err)
	}
}
