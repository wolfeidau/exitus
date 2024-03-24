package jwt

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/wolfeidau/go-oidc"
	jose "gopkg.in/square/go-jose.v2"
)

var supportedSigAlgs = []string{"RS256"}

// Validate returns validates the token, then returns just the parsed JSON from the JWT.
func Validate(ctx context.Context, providerURL, token string) (*JwtPayload, error) {
	// triggers a web request
	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		return nil, errors.Wrap(err, "exitus: failed to get endpoints")
	}

	jws, err := jose.ParseSigned(token)
	if err != nil {
		return nil, errors.Wrap(err, "exitus: failed to validate token")
	}

	// validate signature exits
	switch len(jws.Signatures) {
	case 0:
		return nil, errors.New("exitus: id token not signed")
	case 1:
	default:
		return nil, errors.New("exitus: multiple signatures on id token not supported")
	}

	sig := jws.Signatures[0]
	if !contains(supportedSigAlgs, sig.Header.Algorithm) {
		return nil, errors.Errorf("exitus: id token signed with unsupported algorithm, expected %q got %q", supportedSigAlgs, sig.Header.Algorithm)
	}

	// Throw out tokens with invalid claims before trying to verify the token. This lets
	// us do cheap checks before possibly re-syncing keys.
	payload, err := parseJWT(token)
	if err != nil {
		return nil, errors.Wrap(err, "exitus: failed to validate token")
	}

	jwtp := new(JwtPayload)
	err = json.Unmarshal(payload, jwtp)
	if err != nil {
		return nil, errors.New("exitus: failed to parse jwt payload")
	}

	// Check issuer.
	if jwtp.Issuer != providerURL {
		return nil, errors.Errorf("exitus: failed to match issuer expected: %s actual: %s", providerURL, jwtp.Issuer)
	}

	// Check if the token is expired
	if jwtp.Expires.Time().Before(time.Now()) {
		return nil, errors.Errorf("exitus: token expired current: %s actual: %s", time.Now(), jwtp.Expires.Time())
	}

	// validate signature
	gotPayload, err := provider.VerifySignature(ctx, token)
	if err != nil {
		return nil, errors.Wrap(err, "exitus: failed to validate token")
	}

	// Ensure that the payload returned actually matches the payload parsed earlier.
	if !bytes.Equal(gotPayload, payload) {
		return nil, errors.New("exitus: internal error, payload parsed did not match previous payload")
	}

	return jwtp, nil
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

type JwtPayload struct {
	Sub      string   `json:"sub"`
	TokenUse string   `json:"token_use"`
	Scope    string   `json:"scope"`
	AuthTime JSONTime `json:"auth_time"`
	Issuer   string   `json:"iss"`
	Expires  JSONTime `json:"exp"`
	IssuedAt JSONTime `json:"iat"`
	Version  int      `json:"version"`
	Jti      string   `json:"jti"`
	ClientID string   `json:"client_id"`
}

type JSONTime time.Time

func (j *JSONTime) UnmarshalJSON(b []byte) error {
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
	*j = JSONTime(time.Unix(unix, 0))
	return nil
}

func (j *JSONTime) Time() time.Time {
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

func SplitScopes(scope string) []string {
	return strings.Split(scope, " ")
}
