package handlers

import (
	"context"

	"github.com/iamnoturkkitty/shortener/internal/config"
	"github.com/iamnoturkkitty/shortener/internal/domain/links"
)

type Repository interface {
	SaveLink(ctx context.Context, l links.Link) error
	GetLink(ctx context.Context, hash string) (*links.Link, error)
	GetLinkByUserID(ctx context.Context, userID int) ([]*links.Link, error)
	SaveLinkBatch(ctx context.Context, l []links.Link) error
	DeleteLinkBatch(ctx context.Context, ls []links.DeleteLinkTask) error
	Test() error
}

type Handler struct {
	repo        Repository
	baseAddress string
	deleteQueue chan<- links.DeleteLinkTask
}

func NewHandler(repo Repository, cfg *config.Config, queue chan<- links.DeleteLinkTask) *Handler {
	return &Handler{
		repo:        repo,
		baseAddress: cfg.BaseAddress,
		deleteQueue: queue,
	}
}
