// Copyright 2026 The k8shell CLI Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/config"
)

// Client is an authenticated HTTP client for the k8shell API.
type Client struct {
	server string
	token  string
	debug  bool
	http   *http.Client
}

// New creates a Client from the given context, with optional debug logging and TLS verification skipping.
func New(ctx *config.Context, debug, insecure bool) *Client {
	return &Client{
		server: ctx.Server,
		token:  ctx.Token,
		debug:  debug,
		http:   newHTTPClient(insecure),
	}
}

// NewAnonymous creates a Client for unauthenticated requests (e.g. the login flow).
func NewAnonymous(server string, debug, insecure bool) *Client {
	return &Client{server: server, debug: debug, http: newHTTPClient(insecure)}
}

// newHTTPClient returns a plain http.Client, or one that skips TLS verification when insecure is true.
func newHTTPClient(insecure bool) *http.Client {
	if !insecure {
		return &http.Client{}
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		},
	}
}

// userPath returns the API base path for a user-scoped resource.
// An empty username resolves to /api/v1/me (the token's own user).
func (c *Client) userPath(username string) string {
	if username == "" {
		return "/api/v1/me"
	}
	return "/api/v1/users/" + username
}

// maskToken returns the first six characters of the token followed by "***", for safe logging.
func (c *Client) maskToken() string {
	if len(c.token) <= 6 {
		return "***"
	}
	return c.token[:6] + "***"
}

// debugRequest writes outgoing request headers to stderr.
func (c *Client) debugRequest(req *http.Request) {
	fmt.Fprintf(os.Stderr, "> %s %s\n", req.Method, req.URL)
	keys := make([]string, 0, len(req.Header))
	for k := range req.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := req.Header.Get(k)
		if k == "Authorization" {
			v = "Bearer " + c.maskToken()
		}
		fmt.Fprintf(os.Stderr, "> %s: %s\n", k, v)
	}
	fmt.Fprintln(os.Stderr, ">")
}

// debugResponse writes incoming response headers to stderr.
func (c *Client) debugResponse(resp *http.Response) {
	fmt.Fprintf(os.Stderr, "< %s\n", resp.Status)
	keys := make([]string, 0, len(resp.Header))
	for k := range resp.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(os.Stderr, "< %s: %s\n", k, resp.Header.Get(k))
	}
	fmt.Fprintln(os.Stderr, "<")
}

// APIError is returned for non-2xx responses and carries the HTTP status code and an optional message from the response body.
type APIError struct {
	StatusCode int
	Message    string
}

// Error implements the error interface, mapping common HTTP status codes to human-readable messages.
func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	switch e.StatusCode {
	case http.StatusUnauthorized:
		return "unauthorized — verify your PAT token"
	case http.StatusForbidden:
		return "access denied"
	case http.StatusNotFound:
		return "not found"
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		return fmt.Sprintf("server error (%d)", e.StatusCode)
	default:
		return fmt.Sprintf("request failed (%d)", e.StatusCode)
	}
}

// get performs an authenticated GET request and JSON-decodes the response body into out.
func (c *Client) get(path string, out any) error {
	req, err := http.NewRequest(http.MethodGet, c.server+path, nil)
	if err != nil {
		return err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/json")

	if c.debug {
		c.debugRequest(req)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.debug {
		c.debugResponse(resp)
	}

	if resp.StatusCode >= 400 {
		return &APIError{StatusCode: resp.StatusCode}
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

// post performs an authenticated POST request, JSON-encoding body, and decodes the response into out.
// Pass nil for out to discard the response body.
func (c *Client) post(path string, body any, out any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.server+path, bytes.NewReader(b))
	if err != nil {
		return err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if c.debug {
		c.debugRequest(req)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.debug {
		c.debugResponse(resp)
	}

	if resp.StatusCode >= 400 {
		return &APIError{StatusCode: resp.StatusCode}
	}
	if out != nil && resp.StatusCode != http.StatusNoContent {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

// delete performs an authenticated DELETE request and discards the response body.
func (c *Client) delete(path string) error {
	req, err := http.NewRequest(http.MethodDelete, c.server+path, nil)
	if err != nil {
		return err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	if c.debug {
		c.debugRequest(req)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.debug {
		c.debugResponse(resp)
	}

	if resp.StatusCode >= 400 {
		return &APIError{StatusCode: resp.StatusCode}
	}
	return nil
}

type providerInfo struct {
	Name         string   `json:"name"`
	Capabilities []string `json:"capabilities"`
}

// ListProviders returns the names of providers that support the OnboardUserWebFlow capability.
func (c *Client) ListProviders() ([]string, error) {
	var providers []providerInfo
	if err := c.get("/api/v1/auth/providers", &providers); err != nil {
		return nil, err
	}
	var names []string
	for _, p := range providers {
		for _, cap := range p.Capabilities {
			if cap == "OnboardUserWebFlow" {
				names = append(names, p.Name)
				break
			}
		}
	}
	return names, nil
}

// PollToken checks whether the PAT for the given OAuth state is ready.
// Returns (nil, nil) when still pending (202), or (token, nil) when ready (200).
func (c *Client) PollToken(state string) (*models.UserToken, error) {
	u, err := url.Parse(c.server + "/api/v1/auth/token")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("state", state)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.debug {
		c.debugRequest(req)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if c.debug {
		c.debugResponse(resp)
	}
	if resp.StatusCode == http.StatusAccepted {
		return nil, nil // still pending
	}
	if resp.StatusCode >= 400 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		var body struct {
			Msg string `json:"msg"`
		}
		if json.NewDecoder(resp.Body).Decode(&body) == nil && body.Msg != "" {
			apiErr.Message = body.Msg
		}
		return nil, apiErr
	}
	var token models.UserToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}
