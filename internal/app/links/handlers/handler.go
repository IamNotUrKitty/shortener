package handlers

import (
	"context"

	"github.com/iamnoturkkitty/shortener/internal/config"
	"github.com/iamnoturkkitty/shortener/internal/domain/links"
)

type Repository interface {
	SaveLink(ctx context.Context, l links.Link) error
	GetLink(ctx context.Context, hash string) (*links.Link, error)
	SaveLinkBatch(ctx context.Context, l []links.Link) error
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
