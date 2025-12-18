// Package tempmailchecker provides a Go client for the TempMailChecker API.
// It allows you to detect disposable/temporary email addresses in real-time.
//
// Basic usage:
//
//	checker := tempmailchecker.New("your_api_key")
//	result, err := checker.Check("test@tempmail.com")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if result.Temp {
//	    fmt.Println("Disposable email detected!")
//	}
package tempmailchecker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// API endpoint constants for different regions.
// Choose the endpoint closest to your users for lowest latency.
const (
	// EndpointEU is the European endpoint (default). Best for EU, Africa, Middle East.
	EndpointEU = "https://tempmailchecker.com"
	// EndpointUS is the United States endpoint. Best for Americas.
	EndpointUS = "https://us.tempmailchecker.com"
	// EndpointAsia is the Asia endpoint. Best for Asia-Pacific, Australia, Japan.
	EndpointAsia = "https://asia.tempmailchecker.com"
)

const (
	defaultTimeout = 10 * time.Second
	userAgent      = "TempMailChecker-Go/1.0.0"
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// Client is the TempMailChecker API client.
type Client struct {
	apiKey   string
	endpoint string
	timeout  time.Duration
	client   *http.Client
}

// CheckResult represents the result of an email check.
type CheckResult struct {
	// Temp is true if the email is from a disposable/temporary email provider.
	Temp bool `json:"temp"`
}

// UsageResult represents the API usage statistics.
type UsageResult struct {
	// UsageToday is the number of requests made today.
	UsageToday int `json:"usage_today"`
	// Limit is the daily request limit.
	Limit int `json:"limit"`
	// Reset indicates when the usage counter resets.
	Reset string `json:"reset"`
}

// Option is a function that configures the Client.
type Option func(*Client)

// WithEndpoint sets a custom API endpoint.
func WithEndpoint(endpoint string) Option {
	return func(c *Client) {
		c.endpoint = endpoint
	}
}

// WithTimeout sets a custom timeout for API requests.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
		c.client.Timeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.client = httpClient
	}
}

// New creates a new TempMailChecker client with the given API key.
// Optional configuration can be provided using Option functions.
//
// Example:
//
//	// Basic usage with default settings
//	checker := tempmailchecker.New("your_api_key")
//
//	// With custom endpoint and timeout
//	checker := tempmailchecker.New("your_api_key",
//	    tempmailchecker.WithEndpoint(tempmailchecker.EndpointUS),
//	    tempmailchecker.WithTimeout(5 * time.Second),
//	)
func New(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" || strings.TrimSpace(apiKey) == "" {
		return nil, ErrAPIKeyRequired
	}

	c := &Client{
		apiKey:   apiKey,
		endpoint: EndpointEU,
		timeout:  defaultTimeout,
		client: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// MustNew creates a new TempMailChecker client and panics if the API key is invalid.
// Use this for initialization where you want to fail fast.
func MustNew(apiKey string, opts ...Option) *Client {
	c, err := New(apiKey, opts...)
	if err != nil {
		panic(err)
	}
	return c
}

// Check verifies if an email address is from a disposable email provider.
// Returns a CheckResult with Temp=true if the email is disposable.
func (c *Client) Check(email string) (*CheckResult, error) {
	if email == "" || strings.TrimSpace(email) == "" {
		return nil, ErrEmailRequired
	}

	email = strings.TrimSpace(email)
	if !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmail
	}

	apiURL := fmt.Sprintf("%s/check?email=%s", c.endpoint, url.QueryEscape(email))
	return c.doCheckRequest(apiURL)
}

// CheckDomain verifies if a domain is a disposable email provider.
// Returns a CheckResult with Temp=true if the domain is disposable.
func (c *Client) CheckDomain(domain string) (*CheckResult, error) {
	if domain == "" || strings.TrimSpace(domain) == "" {
		return nil, ErrDomainRequired
	}

	domain = strings.TrimSpace(domain)
	apiURL := fmt.Sprintf("%s/check?domain=%s", c.endpoint, url.QueryEscape(domain))
	return c.doCheckRequest(apiURL)
}

// IsDisposable is a convenience method that returns true if the email is disposable.
func (c *Client) IsDisposable(email string) (bool, error) {
	result, err := c.Check(email)
	if err != nil {
		return false, err
	}
	return result.Temp, nil
}

// GetUsage retrieves the current API usage statistics.
func (c *Client) GetUsage() (*UsageResult, error) {
	apiURL := fmt.Sprintf("%s/usage?key=%s", c.endpoint, url.QueryEscape(c.apiKey))

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp.StatusCode, body)
	}

	var result UsageResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

func (c *Client) doCheckRequest(apiURL string) (*CheckResult, error) {
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp.StatusCode, body)
	}

	var result CheckResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

func (c *Client) parseError(statusCode int, body []byte) error {
	var errResp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &errResp); err == nil {
		if statusCode == http.StatusTooManyRequests {
			msg := errResp.Message
			if msg == "" {
				msg = "Daily limit reached"
			}
			return &RateLimitError{Message: msg}
		}

		if errResp.Error != "" {
			return &APIError{
				StatusCode: statusCode,
				Message:    errResp.Error,
			}
		}
	}

	return &APIError{
		StatusCode: statusCode,
		Message:    fmt.Sprintf("API request failed with status %d", statusCode),
	}
}

