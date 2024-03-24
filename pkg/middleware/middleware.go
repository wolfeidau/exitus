package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

// RequestID middleware which checks request id.
func RequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		rid := req.Header.Get("X-Amzn-Trace-Id")

		if rid == "" {
			rid = uuid.NewV4().String()
		}

		res.Header().Set("X-Amzn-Trace-Id", rid)

		return next(c)
	}
}

// ErrorLog middleware which logs each error.
func ErrorLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		if err = next(c); err != nil {
			log.Info().Str("Trace-Id", c.Response().Header().Get("X-Amzn-Trace-Id")).Msgf("error occurred: %+v", err)
			c.Error(err)
		}

		return
	}
}

// RequestLog middleware which logs each request.
func RequestLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		start := time.Now()

		if err = next(c); err != nil {
			c.Error(err)
		}

		stop := time.Now()

		traceid := c.Response().Header().Get("X-Amzn-Trace-Id")
		path := c.Request().URL.Path
		method := c.Request().Method

		log.Info().Str("X-Amzn-Trace-Id", traceid).Str("RealIP", c.RealIP()).Str("method", method).Str("path", path).Int("status", c.Response().Status).Int64("length", c.Response().Size).Str("latency", stop.Sub(start).String()).Msg("request")

		return
	}
}
