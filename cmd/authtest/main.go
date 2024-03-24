package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/motemen/go-loghttp"
	"github.com/pkg/errors"
	"github.com/wolfeidau/go-oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	jose "gopkg.in/square/go-jose.v2"
)

var supportedSigAlgs = []string{"RS256"}

func main() {
	appClientID := os.Getenv("OAUTH_CLIENT_ID")
	appClientSecret := os.Getenv("OAUTH_CLIENT_SECRET")

	providerURL := os.Getenv("OPENID_PROVIDER_URL")

	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		log.Fatalf("failed to get endpoints: %+v", err)
	}

	conf := &clientcredentials.Config{
		ClientID:     appClientID,
		ClientSecret: appClientSecret,
		Scopes:       []string{"exitus/issue.read", "exitus/issue.write", "exitus/project.read", "exitus/project.write"},
		TokenURL:     provider.Endpoint().TokenURL,
		AuthStyle:    oauth2.AuthStyleInHeader,
	}

	http.DefaultTransport = &loghttp.Transport{
		Transport: http.DefaultTransport,
		LogRequest: func(req *http.Request) {
			// log.Printf("[%p] %s %s", req, req.Method, req.URL)
			data, _ := httputil.DumpRequest(req, true)
			log.Println(string(data))
		},
		LogResponse: func(resp *http.Response) {
			// log.Printf("[%p] %d %s", resp.Request, resp.StatusCode, resp.Request.URL)
		},
	}

	t, err := conf.Token(ctx)
	if err != nil {
		log.Fatalf("failed to get token: %+v", err)
	}

	log.Println(t.AccessToken)

	jws, err := jose.ParseSigned(t.AccessToken)
	if err != nil {
		log.Fatalf("failed to validate token: %+v", err)
	}

	// validate signature exits
	switch len(jws.Signatures) {
	case 0:
		log.Fatalf("exitus: id token not signed")
	case 1:
	default:
		log.Fatalf("exitus: multiple signatures on id token not supported")
	}

	sig := jws.Signatures[0]
	if !contains(supportedSigAlgs, sig.Header.Algorithm) {
		log.Fatalf("exitus: id token signed with unsupported algorithm, expected %q got %q", supportedSigAlgs, sig.Header.Algorithm)
	}

	// Throw out tokens with invalid claims before trying to verify the token. This lets
	// us do cheap checks before possibly re-syncing keys.
	payload, err := parseJWT(t.AccessToken)
	if err != nil {
		log.Fatalf("failed to validate token: %+v", err)
	}

	jwtp := new(jwtPayload)
	err = json.Unmarshal(payload, jwtp)
	if err != nil {
		log.Fatal("exitus: failed to parse jwt payload")
	}

	// Check issuer.
	if jwtp.Issuer != providerURL {
		log.Fatalf("failed to match issuer expected: %s actual: %s", providerURL, jwtp.Issuer)
	}

	// Check if the token is expired
	if jwtp.Expires.Time().Before(time.Now()) {
		log.Fatalf("token expired current: %s actual: %s", time.Now(), jwtp.Expires.Time())
	}

	// validate signature
	gotPayload, err := provider.VerifySignature(ctx, t.AccessToken)
	if err != nil {
		log.Fatalf("failed to validate token: %+v", err)
	}

	// Ensure that the payload returned actually matches the payload parsed earlier.
	if !bytes.Equal(gotPayload, payload) {
		log.Fatalf("exitus: internal error, payload parsed did not match previous payload")
	}

	log.Println(jwtp)
}

func parseJWT(p string) ([]byte, error) {
	parts := strings.Split(p, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("exitus: malformed jwt, expected 3 parts got %d", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.Wrapf(err, "exitus: malformed jwt payload: %v", err)
	}

	return payload, nil
}

type jwtPayload struct {
	Sub      string   `json:"sub"`
	TokenUse string   `json:"token_use"`
	Scope    string   `json:"scope"`
	AuthTime jsonTime `json:"auth_time"`
	Issuer   string   `json:"iss"`
	Expires  jsonTime `json:"exp"`
	IssuedAt jsonTime `json:"iat"`
	Version  int      `json:"version"`
	Jti      string   `json:"jti"`
	ClientID string   `json:"client_id"`
}

type jsonTime time.Time

func (j *jsonTime) UnmarshalJSON(b []byte) error {
	var n json.Number
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	var unix int64

	if t, err := n.Int64(); err == nil {
		unix = t
	} else {
		f, err := n.Float64()
		if err != nil {
			return err
		}
		unix = int64(f)
	}
	*j = jsonTime(time.Unix(unix, 0))
	return nil
}

func (j *jsonTime) Time() time.Time {
	return time.Time(*j)
}

func contains(sli []string, ele string) bool {
	for _, s := range sli {
		if s == ele {
			return true
		}
	}
	return false
}
