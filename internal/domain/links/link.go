package links

import (
	"errors"
	"net/url"

	"github.com/sqids/sqids-go"
)

var (
	ErrLinkNotFound = errors.New("не найден URL")
	ErrBadURL       = errors.New("некорректный URL")
	ErrLinkCreation = errors.New("ошибка создания короткой ссылки")
)

var s, _ = sqids.New()

func makeHash(byteURL []byte) (string, error) {
	d := make([]uint64, len(byteURL))

	for i, b := range byteURL {
		d[i] = uint64(b)
	}

	hash, err := s.Encode(d)

	return hash[:6], err
}

func validateURL(urlString string) error {
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return ErrBadURL
	}

	return nil
}

type Link struct {
	hash string
	url  string
}

func NewLink(url string) (*Link, error) {
	errURL := validateURL(url)
	if errURL != nil {
		return nil, errURL
	}

	hash, err := makeHash([]byte(url))
	if err != nil {
		return nil, ErrLinkCreation
	}

	return &Link{
		hash: hash,
		url:  url,
	}, nil
}

func (l *Link) Hash() string {
	return l.hash
}

func (l *Link) URL() string {
	return l.url
}
