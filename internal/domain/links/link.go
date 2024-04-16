package links

import (
	"errors"
	"net/url"

	"github.com/sqids/sqids-go"
)

var (
	LinkNotFound = errors.New("link not found")
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
		return errors.New("Bad url")
	}

	return nil
}

type Link struct {
	hash string
	url  string
}

func NewLink(url string) (*Link, error) {
	hash, err := makeHash([]byte(url))
	if err != nil {
		return nil, err
	}

	errURL := validateURL(url)
	if errURL != nil {
		return nil, errURL
	}

	return &Link{
		hash: hash,
		url:  url,
	}, nil
}

func (l *Link) Hash() string {
	return l.hash
}

func (l *Link) Url() string {
	return l.url
}