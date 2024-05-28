package echomiddleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func InitLoggerMiddleware(l *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:       true,
		LogURI:          true,
		LogResponseSize: true,
		LogMethod:       true,
		LogLatency:      true,
		LogHeaders:      []string{echo.HeaderContentType, echo.HeaderCookie},

		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			l.Info("request",
				zap.String("URI", v.URI),
				zap.String("Method", v.Method),
				zap.Duration("Latency", v.Latency),
				zap.String("Content-type", strings.Join(v.Headers[echo.HeaderContentType], " ")),
				zap.String("Cookie", strings.Join(v.Headers[echo.HeaderCookie], " ")),
			)

			l.Info("response",
				zap.Int("Status", v.Status),
				zap.Int64("Size", v.ResponseSize),
			)
			return nil
		},
	})
}
