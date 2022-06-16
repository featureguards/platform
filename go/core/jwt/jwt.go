package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"time"

	"platform/go/core/env"
	"platform/go/core/ids"
	"platform/go/core/random"

	"github.com/benbjohnson/clock"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type TokenType string

var (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

const (
	// This is for refresh rotation.
	TokenFamilyClaim = "family"
)

type JWT struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	clock      clock.Clock
}

type Options func(j *JWT) error

func WithKeyPairFiles(private, public string) Options {
	return func(j *JWT) error {
		privateContent, err := ioutil.ReadFile(private)
		if err != nil {
			return errors.WithStack(err)
		}
		publicContent, err := ioutil.ReadFile(public)
		if err != nil {
			return errors.WithStack(err)
		}
		privatePEM, _ := pem.Decode(privateContent)
		if err != nil {
			return errors.WithStack(err)
		}

		privateKey, err := x509.ParsePKCS1PrivateKey(privatePEM.Bytes)
		if err != nil {
			return errors.WithStack(err)
		}

		publicPEM, _ := pem.Decode(publicContent)
		if err != nil {
			return errors.WithStack(err)
		}
		publicKey, err := x509.ParsePKIXPublicKey(publicPEM.Bytes)
		if err != nil {
			return errors.WithStack(err)
		}
		j.privateKey = privateKey
		j.publicKey = publicKey.(*rsa.PublicKey)
		return nil
	}
}

func WithKeyPair(private *rsa.PrivateKey, public *rsa.PublicKey) Options {
	return func(j *JWT) error {
		j.privateKey = private
		j.publicKey = public
		return nil
	}
}

func WithClock(c clock.Clock) Options {
	return func(j *JWT) error {
		j.clock = c
		return nil
	}
}

func New(options ...Options) (*JWT, error) {
	j := &JWT{}
	for _, opt := range options {
		if err := opt(j); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if j.privateKey == nil || j.publicKey == nil {
		return nil, errors.WithStack(fmt.Errorf("no private/public key"))
	}
	return j, nil
}

type TokenOptions func(to *tokenOptions) error
type tokenOptions struct {
	family string
}

func WithFamily(family string) TokenOptions {
	return func(to *tokenOptions) error {
		to.family = family
		return nil
	}
}

func (j *JWT) SignedToken(apiKey ids.ID, tokenType TokenType, options ...TokenOptions) ([]byte, error) {
	opts := &tokenOptions{}
	for _, opt := range options {
		opt(opts)
	}
	if opts.family == "" {
		opts.family = random.RandString(16, nil)
	}
	t, err := jwt.NewBuilder().Issuer(env.Domain).IssuedAt(j.clock.Now().UTC()).
		NotBefore(j.clock.Now().UTC()).Subject(string(apiKey)).
		JwtID(random.RandString(16, nil)).
		Audience([]string{string(tokenType)}).Expiration(exp(j.clock, tokenType)).Build()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := t.Set(TokenFamilyClaim, opts.family); err != nil {
		return nil, errors.WithStack(err)
	}
	signed, err := jwt.Sign(t, jwt.WithKey(jwa.RS256, j.privateKey))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return signed, nil
}

func (j *JWT) ParseToken(payload []byte) (jwt.Token, error) {
	token, err := jwt.Parse(
		payload,
		jwt.WithValidate(true),
		jwt.WithClock(j.clock),
		jwt.WithAcceptableSkew(10*time.Second),
		jwt.WithKey(jwa.RS256, j.publicKey),
		jwt.WithIssuer(env.Domain),
	)
	if err != nil {
		log.Warnf("Invalid JWT token: %s\n", err)
		return nil, errors.WithStack(err)
	}
	return token, nil
}

func exp(c clock.Clock, t TokenType) time.Time {
	switch t {
	case RefreshToken:
		return c.Now().UTC().Add(7 * 24 * time.Hour)
	default:
		return c.Now().UTC().Add(15 * time.Minute)
	}
}
