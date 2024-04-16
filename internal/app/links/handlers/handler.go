package handlers

import "github.com/iamnoturkkitty/shortener/internal/domain/links"

type Repository interface {
	SaveLink(l links.Link) error
	GetLink(hash string) (*links.Link, error)
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}
