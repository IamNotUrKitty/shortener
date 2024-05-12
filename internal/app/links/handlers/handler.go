package handlers

import (
	"github.com/iamnoturkkitty/shortener/internal/config"
	"github.com/iamnoturkkitty/shortener/internal/domain/links"
)

type Repository interface {
	SaveLink(l links.Link) error
	GetLink(hash string) (*links.Link, error)
	Test() error
}

type Handler struct {
	repo        Repository
	baseAddress string
}

func NewHandler(repo Repository, cfg *config.Config) *Handler {
	return &Handler{
		repo:        repo,
		baseAddress: cfg.BaseAddress,
	}
}
