package internal

import (
	"strings"

	"github.com/iamnoturkkitty/shortener/internal/app/links"
	"github.com/iamnoturkkitty/shortener/internal/config"
	linksInfra "github.com/iamnoturkkitty/shortener/internal/infrastructure/links"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func NewServer(cfg *config.Config) *echo.Echo {
	e := echo.New()

	logger, _ := zap.NewDevelopment()

	//TODO: move logger init to some utils
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:       true,
		LogURI:          true,
		LogResponseSize: true,
		LogMethod:       true,
		LogLatency:      true,
		LogHeaders:      []string{echo.HeaderContentType},

		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.String("Method", v.Method),
				zap.Duration("Latency", v.Latency),
				zap.String("Content-type", strings.Join(v.Headers[echo.HeaderContentType], " ")),
			)

			logger.Info("response",
				zap.Int("Status", v.Status),
				zap.Int64("Size", v.ResponseSize),
			)
			return nil
		},
	}))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			h := c.Request().Header.Get(echo.HeaderContentType)
			return h != echo.MIMEApplicationJSON && h != echo.MIMETextHTML
		},
	}))

	linksRepo := linksInfra.NewInMemoryRepo()

	links.Setup(e, linksRepo, cfg)

	return e
}
