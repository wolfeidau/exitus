package auth

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

const (
	// UserKey the key used in the context for the user.
	UserKey = "Auth.User"
)

var (
	// ErrUserNotFound user not found in context.
	ErrUserNotFound = errors.New("user not found in context")

	// ErrUserTypeMismatch failed to retrieve user, type mismatch.
	ErrUserTypeMismatch = errors.New("failed to retrieve user, type mismatch")
)

// AuthenticatedUser authenticated user used to validate access.
type AuthenticatedUser struct {
	ID     string   `json:"id,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
}

// MarshalZerologObject used to print user in logs.
func (au *AuthenticatedUser) MarshalZerologObject(e *zerolog.Event) {
	e.Str("id", au.ID).Strs("scopes", au.Scopes)
}

// HasScope check if the authenticated user has one of the allowed
// scopes.
func (au *AuthenticatedUser) HasScope(allowedScopes []string) bool {
	if len(allowedScopes) == 0 {
		return false
	}

	// O(n^2)
	for _, usc := range au.Scopes {
		for _, asc := range allowedScopes {
			if usc == asc {
				return true
			}
		}
	}

	return false
}

// LoadUserFromContext load the authenticated user from the context.
func LoadUserFromContext(ctx echo.Context) (AuthenticatedUser, error) {
	var (
		user AuthenticatedUser
		ok   bool
	)

	uval := ctx.Get(UserKey)
	if uval == nil {
		return user, ErrUserNotFound
	}

	user, ok = uval.(AuthenticatedUser)
	if !ok {
		return user, ErrUserTypeMismatch
	}

	return user, nil
}
