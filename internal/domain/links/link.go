package links

import (
	"errors"
	"net/url"

	"github.com/google/uuid"
	"github.com/sqids/sqids-go"
)

var (
	ErrLinkNotFound  = errors.New("не найден URL")
	ErrBadURL        = errors.New("некорректный URL")
	ErrLinkCreation  = errors.New("ошибка создания короткой ссылки")
	ErrLinkDuplicate = errors.New("короткая ссылка уже была создана")
)

var s, _ = sqids.New()

func makeHash(byteURL []byte) (string, error) {
	d := make([]uint64, len(byteURL))

	for i, b := range byteURL {
		d[i] = uint64(b)
	}

	hash, err := s.Encode(d)

	return hash[len(hash)-6:], err
}

func validateURL(urlString string) error {
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return ErrBadURL
	}

	return nil
}

type StoredLink struct {
	Hash string    `json:"short_url"`
	URL  string    `json:"original_url"`
	ID   uuid.UUID `json:"uuid"`
}

type Link struct {
	hash string
	url  string
	id   uuid.UUID
}

func NewLink(id uuid.UUID, url, hash string) (*Link, error) {
	return &Link{
		id:   id,
		hash: hash,
		url:  url,
	}, nil
}

func CreateLink(url string) (*Link, error) {
	if err := validateURL(url); err != nil {
		return nil, err
	}

	hash, err := makeHash([]byte(url))
	if err != nil {
		return nil, ErrLinkCreation
	}

	return NewLink(uuid.New(), url, hash)
}

func (l *Link) Hash() string {
	return l.hash
}

func (l *Link) URL() string {
	return l.url
}

func (l *Link) ID() uuid.UUID {
	return l.id
}
