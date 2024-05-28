package internal

import (
	"net/http"

	"github.com/iamnoturkkitty/shortener/internal/app/links"
	"github.com/iamnoturkkitty/shortener/internal/config"
	linksDomain "github.com/iamnoturkkitty/shortener/internal/domain/links"
	echomiddleware "github.com/iamnoturkkitty/shortener/internal/echomiddleware"
	linksInfra "github.com/iamnoturkkitty/shortener/internal/infrastructure/links"
	"github.com/iamnoturkkitty/shortener/internal/workers"
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

	e.Use(echomiddleware.InitJWTMiddleware())

	linksRepo, err := linksInfra.Setup(cfg)
	if err != nil {
		return nil, err
	}

	deleteQueue := make(chan linksDomain.DeleteLinkTask)

	links.Setup(e, linksRepo, cfg, deleteQueue)

	go workers.DeleteLinksWorker(deleteQueue, linksRepo)

	e.GET("/ping", func(c echo.Context) error {
		if err := linksRepo.Test(); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, "ok")
	})

	return e, nil
}
