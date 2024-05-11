package internal

import (
	"github.com/iamnoturkkitty/shortener/internal/app/links"
	"github.com/iamnoturkkitty/shortener/internal/config"
	echomiddleware "github.com/iamnoturkkitty/shortener/internal/echomiddleware"
	linksInfra "github.com/iamnoturkkitty/shortener/internal/infrastructure/links"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func NewServer(cfg *config.Config) (*echo.Echo, error) {
	e := echo.New()

	logger, err := zap.NewDevelopment()

	if err != nil {
		return nil, err
	}

	e.Use(echomiddleware.InitLoggerMiddleware(logger))

	e.Use(echomiddleware.InitGzipMiddleware())

	e.Use(middleware.Decompress())

	linksRepo, err := linksInfra.InitFSRepo(cfg.StorageFile)
	if err != nil {
		return nil, err
	}

	links.Setup(e, linksRepo, cfg)

	return e, nil
}
