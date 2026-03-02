package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/voska/amadeus-cli/internal/config"
)

type Token struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int       `json:"expires_in"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func FetchToken(apiKey, apiSecret, env string) (*Token, error) {
	base := config.BaseURL(env)
	endpoint := base + "/v1/security/oauth2/token"

	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {apiKey},
		"client_secret": {apiSecret},
	}

	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request returned %d: %s", resp.StatusCode, string(body))
	}

	var tok Token
	if err := json.Unmarshal(body, &tok); err != nil {
		return nil, fmt.Errorf("parsing token response: %w", err)
	}
	tok.ExpiresAt = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
	return &tok, nil
}
