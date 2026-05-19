package oauth2

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type GrantType string
type TokenType string

const (
	GrantAuthCode    GrantType = "authorization_code"
	GrantClientCreds GrantType = "client_credentials"
	TokenBearer      TokenType = "Bearer"
)

type AuthCode struct {
	Code        string    `json:"code"`
	ClientID    string    `json:"client_id"`
	RedirectURI string    `json:"redirect_uri"`
	Scope       []string  `json:"scope"`
	Subject     string    `json:"subject"`
	ExpiresAt   time.Time `json:"expires_at"`
	Used        bool      `json:"used"`
}

type TokenIntrospection struct {
	Active   bool     `json:"active"`
	Subject  string   `json:"sub"`
	ClientID string   `json:"client_id"`
	Scope    []string `json:"scope"`
	Exp      int64    `json:"exp"`
	TokenType TokenType `json:"token_type"`
}

type AccessToken struct {
	Token     string    `json:"access_token"`
	TokenType TokenType `json:"token_type"`
	ExpiresIn int       `json:"expires_in"`
	Scope     []string  `json:"scope"`
	IssuedAt  time.Time `json:"issued_at"`
}

func IssueCode(clientID, redirectURI, subject string, scope []string) AuthCode {
	return AuthCode{
		Code:        uuid.NewString(),
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		Subject:     subject,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}
}

func ExchangeCode(code *AuthCode, clientID string) (AccessToken, error) {
	if code.Used {
		return AccessToken{}, errors.New("auth code already used")
	}
	if code.ClientID != clientID {
		return AccessToken{}, errors.New("client_id mismatch")
	}
	if time.Now().After(code.ExpiresAt) {
		return AccessToken{}, errors.New("auth code expired")
	}
	code.Used = true
	return AccessToken{
		Token:     uuid.NewString(),
		TokenType: TokenBearer,
		ExpiresIn: 3600,
		Scope:     code.Scope,
		IssuedAt:  time.Now(),
	}, nil
}

func Introspect(rawToken, expectedClientID string, knownTokens map[string]AccessToken) TokenIntrospection {
	if rawToken == "" {
		return TokenIntrospection{Active: false}
	}
	rawToken = strings.TrimPrefix(rawToken, "Bearer ")
	tok, ok := knownTokens[rawToken]
	if !ok {
		return TokenIntrospection{Active: false}
	}
	if time.Now().After(tok.IssuedAt.Add(time.Duration(tok.ExpiresIn) * time.Second)) {
		return TokenIntrospection{Active: false}
	}
	return TokenIntrospection{
		Active:    true,
		ClientID:  expectedClientID,
		Scope:     tok.Scope,
		Exp:       tok.IssuedAt.Add(time.Duration(tok.ExpiresIn) * time.Second).Unix(),
		TokenType: tok.TokenType,
	}
}
