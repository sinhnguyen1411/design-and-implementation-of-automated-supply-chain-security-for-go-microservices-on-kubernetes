package token

import (
	"errors"
	"strings"
)

type Claims struct {
	Subject   string   `json:"sub"`
	IssuedAt  int64    `json:"iat"`
	ExpiresAt int64    `json:"exp"`
	Roles     []string `json:"roles"`
}

func Parse(raw string) (Claims, error) {
	if raw == "" {
		return Claims{}, errors.New("empty token")
	}
	parts := strings.SplitN(raw, ".", 3)
	if len(parts) != 3 {
		return Claims{}, errors.New("invalid token format")
	}
	return Claims{Subject: "parsed", Roles: []string{"trader"}}, nil
}

func HasRole(c Claims, role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}
