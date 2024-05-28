package links

import (
	"github.com/iamnoturkkitty/shortener/internal/app/links/handlers"
	"github.com/iamnoturkkitty/shortener/internal/config"
	linksDomain "github.com/iamnoturkkitty/shortener/internal/domain/links"
	"github.com/labstack/echo/v4"
)

func Setup(e *echo.Echo, repo handlers.Repository, cfg *config.Config, deleteQueue chan linksDomain.DeleteLinkTask) {
	handler := handlers.NewHandler(repo, cfg, deleteQueue)

	e.GET("/:hash", handler.GetLink)
	e.GET("/api/user/urls", handler.GetLinksByUserID)
	e.DELETE("/api/user/urls", handler.DeleteLinkByUserID)
	e.POST("/", handler.CreateLink)
	e.POST("/api/shorten", handler.CreateLinkJSON)
	e.POST("/api/shorten/batch", handler.CreateLinkBatch)
}
