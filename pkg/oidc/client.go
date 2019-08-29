package oidc

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/wolfeidau/go-oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Session struct {
	config *clientcredentials.Config
}

type SessionConfig struct {
	ProviderURL  string
	ClientID     string
	ClientSecret string
	Scopes       []string
}

func New(ctx context.Context, config SessionConfig) (*Session, error) {

	provider, err := oidc.NewProvider(ctx, config.ProviderURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get endpoint")
	}

	conf := &clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       config.Scopes,
		TokenURL:     provider.Endpoint().TokenURL,
		AuthStyle:    oauth2.AuthStyleInHeader,
	}

	return &Session{config: conf}, nil
}

func (sess *Session) Client(ctx context.Context) *http.Client {
	return sess.config.Client(ctx)
}
