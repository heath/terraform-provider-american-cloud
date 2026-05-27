package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "https://api.americancloud.com"
	userAgent      = "terraform-provider-americancloud"
)

// Client is an HTTP client for the American Cloud API.
type Client struct {
	BaseURL      string
	BearerToken  string
	ClientID     string
	ClientSecret string
	HTTPClient   *http.Client
}

// NewClient creates a new American Cloud API client.
// It supports two auth modes:
//   - Bearer token (set BearerToken)
//   - API Key (set ClientID + ClientSecret)
func NewClient(baseURL, bearerToken, clientID, clientSecret string) (*Client, error) {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	// Validate the base URL
	_, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL %q: %w", baseURL, err)
	}

	if bearerToken == "" && (clientID == "" || clientSecret == "") {
		return nil, fmt.Errorf("either bearer_token or both client_id and client_secret must be provided")
	}

	return &Client{
		BaseURL:      baseURL,
		BearerToken:  bearerToken,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// doRequest performs an HTTP request and decodes the JSON response.
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	u := c.BaseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reqBody)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Auth: prefer bearer token, fall back to client-id/client-secret
	if c.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.BearerToken)
	} else {
		req.Header.Set("client-id", c.ClientID)
		req.Header.Set("client-secret", c.ClientSecret)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseAPIError(resp.StatusCode, respBody)
	}

	// Some endpoints return 204 No Content
	if resp.StatusCode == 204 || len(respBody) == 0 {
		return nil
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("decoding response (status %d): %w\nbody: %s", resp.StatusCode, err, string(respBody))
		}
	}

	return nil
}

// doRequestWithQuery performs an HTTP request with query parameters.
func (c *Client) doRequestWithQuery(ctx context.Context, method, path string, query url.Values, body interface{}, result interface{}) error {
	if len(query) > 0 {
		path = path + "?" + query.Encode()
	}
	return c.doRequest(ctx, method, path, body, result)
}
