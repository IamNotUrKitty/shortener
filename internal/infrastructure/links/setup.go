package links

import (
	"github.com/iamnoturkkitty/shortener/internal/app/links/handlers"
	"github.com/iamnoturkkitty/shortener/internal/config"
)

func Setup(cfg *config.Config) (handlers.Repository, error) {
	if cfg.DatabaseAddress != "" {
		return NewPostgresRepo(cfg.DatabaseAddress)
	}

	if cfg.StorageFile != "" {
		return NewInFSRepo(cfg.StorageFile)
	}

	return NewInMemoryRepo(), nil
}
