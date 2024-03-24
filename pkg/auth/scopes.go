package auth

import (
	"errors"

	"github.com/labstack/echo/v4"
)

const (
	// OpenIDScopes the key used in the context for openid scopes.
	OpenIDScopes = "OpenId.Scopes"
)

var (
	// ErrScopesNotFound scopes for this operation not found in context.
	ErrScopesNotFound = errors.New("scopes not found in context")
	// ErrScopesTypeMismatch failed to retrieve scopes, type mismatch.
	ErrScopesTypeMismatch = errors.New("failed to retrieve scopes, type mismatch")
)

// LoadOperationScopesFromContext load the scopes for a given operation from the context.
func LoadOperationScopesFromContext(ctx echo.Context, provider string) ([]string, error) {
	var (
		scopes []string
		ok     bool
	)

	sval := ctx.Get(provider)
	if sval == nil {
		return nil, ErrScopesNotFound
	}

	scopes, ok = sval.([]string)
	if !ok {
		return nil, ErrScopesTypeMismatch
	}

	return scopes, nil
}
