package internal

import (
	"github.com/iamnoturkkitty/shortener/internal/app/links"
	"github.com/iamnoturkkitty/shortener/internal/config"
	linksInfra "github.com/iamnoturkkitty/shortener/internal/infrastructure/links"
	"github.com/labstack/echo"
)

func NewServer(cfg config.Config) *echo.Echo {
	e := echo.New()

	linksRepo := linksInfra.NewInMemoryRepo()

	links.Setup(e, linksRepo)

	return e
}
