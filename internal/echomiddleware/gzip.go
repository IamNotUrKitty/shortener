package echomiddleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitGzipMiddleware() echo.MiddlewareFunc {
	return middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			h := c.Request().Header.Get(echo.HeaderContentType)
			return h != echo.MIMEApplicationJSON && h != echo.MIMETextHTML
		},
	})
}
