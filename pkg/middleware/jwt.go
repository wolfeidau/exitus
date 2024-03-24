package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/exitus/pkg/auth"
	"github.com/wolfeidau/exitus/pkg/jwt"
)

const (
	// DefaultAuthScheme default authentication scheme for JWT tokens.
	DefaultAuthScheme = "Bearer"
	// DefaultAuthHeaderName default header to load the JWT.
	DefaultAuthHeaderName = "Authorization"

	// DefaultContextVar the variable in the context to store the JWT token after successful login.
	DefaultContextVar = "user"
)

var (
	// ErrJWTMissing missing or malformed jwt.
	ErrJWTMissing = echo.NewHTTPError(http.StatusBadRequest, "missing or malformed jwt")

	// ErrJWTValidation invalid jwt.
	ErrJWTValidation = echo.NewHTTPError(http.StatusUnauthorized, "invalid jwt")
)

// JWTConfig jwt middleware configuration.
type JWTConfig struct {
	ProviderURL string
	ClientID    string
	AuthScheme  string
}

// JWTWithConfig middleware which validates tokens.
func JWTWithConfig(config *JWTConfig) echo.MiddlewareFunc {
	if config.ProviderURL == "" {
		log.Fatal().Msg("exitus: missing provider URL")
	}
	if config.ClientID == "" {
		log.Fatal().Msg("exitus: missing client ID")
	}
	if config.AuthScheme == "" {
		config.AuthScheme = DefaultAuthScheme
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := extractFromHeader(c, "Authorization", config.AuthScheme)
			if err != nil {
				return err
			}

			jwtp, err := validateToken(c.Request().Context(), config.ProviderURL, token)
			if err != nil {
				return err
			}

			usr := auth.AuthenticatedUser{
				ID:     jwtp.Sub,
				Scopes: jwt.SplitScopes(jwtp.Scope),
			}

			c.Set(auth.UserKey, usr)

			log.Info().Object("user", &usr).Msg("context updated")

			return next(c)
		}
	}
}

// ExtractFromHeader attempt to get the JWT from the provided header.
func extractFromHeader(c echo.Context, header, authScheme string) (string, error) {
	auth := c.Request().Header.Get(header)
	l := len(authScheme)
	if len(auth) > l+1 && auth[:l] == authScheme {
		return auth[l+1:], nil
	}
	return "", ErrJWTMissing
}

// ValidateToken and return an error if it fails.
func validateToken(ctx context.Context, providerURL, token string) (*jwt.JwtPayload, error) {
	payload, err := jwt.Validate(ctx, providerURL, token)
	if err != nil {
		log.Warn().Err(err).Msg("exitus: failed to validate header")
		return nil, ErrJWTValidation
	}
	return payload, nil
}
