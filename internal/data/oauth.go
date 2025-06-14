package data

import (
	"crypto/rand"
	"encoding/base64"
)

type Claim struct {
	Scope       string   `json:"scope"`
	Permissions []string `json:"permissions"`
	UserID      string   `json:"user_id"`
}

type OAuthModel struct {
}

func (m OAuthModel) GenerateState() (string, error) {
	// Generate a random state string for OAuth flow
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(randomBytes), nil
}
