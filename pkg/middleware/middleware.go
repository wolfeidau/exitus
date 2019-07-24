package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// RequestLog middleware adds a `Server` header to the response.
func RequestLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {

		start := time.Now()

		if err = next(c); err != nil {
			c.Error(err)
		}

		stop := time.Now()

		log.Info().Str("RealIP", c.RealIP()).Str("method", c.Request().Method).Str("path", c.Request().URL.Path).Int("status", c.Response().Status).Int64("length", c.Response().Size).Str("latency", stop.Sub(start).String()).Msg("request")

		return
	}
}
