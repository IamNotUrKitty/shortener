package links

import (
	"github.com/iamnoturkkitty/shortener/internal/app/links/handlers"
	"github.com/labstack/echo"
)

func Setup(e *echo.Echo, repo handlers.Repository) {
	handler := handlers.NewHandler(repo)

	e.GET("/:hash", handler.GetLink)
	e.POST("/", handler.CreateLink)
}
