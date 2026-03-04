package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/voska/amadeus-cli/internal/auth"
	"github.com/voska/amadeus-cli/internal/config"
	"github.com/voska/amadeus-cli/internal/errfmt"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	token      *auth.Token
}

func NewClient(token *auth.Token, env string) *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    config.BaseURL(env),
		token:      token,
	}
}

func (c *Client) Get(path string, params url.Values) (map[string]any, error) {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) Post(path string, body any) (map[string]any, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.baseURL+path, strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/vnd.amadeus+json")
	return c.do(req)
}

func (c *Client) do(req *http.Request) (map[string]any, error) {
	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)
	req.Header.Set("Accept", "application/vnd.amadeus+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return result, nil
	case http.StatusUnauthorized:
		return nil, errfmt.Auth("authentication failed — run 'amadeus auth login'")
	case http.StatusForbidden:
		return nil, errfmt.Forbidden(extractError(result))
	case http.StatusNotFound:
		return nil, errfmt.NotFound(extractError(result))
	case http.StatusTooManyRequests:
		return nil, errfmt.RateLimit()
	default:
		if resp.StatusCode >= 500 {
			return nil, errfmt.New(errfmt.ExitRetryable, extractError(result))
		}
		return nil, errfmt.New(errfmt.ExitError, extractError(result))
	}
}

func extractError(result map[string]any) string {
	if errors, ok := result["errors"].([]any); ok && len(errors) > 0 {
		if first, ok := errors[0].(map[string]any); ok {
			if detail, ok := first["detail"].(string); ok {
				return detail
			}
		}
	}
	return "unknown error"
}
