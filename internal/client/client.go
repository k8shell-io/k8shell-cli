package client

import (
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

type Client struct {
	server string
	token  string
	debug  bool
	http   *http.Client
}

func New(ctx *config.Context, debug, insecure bool) *Client {
	return &Client{
		server: ctx.Server,
		token:  ctx.Token,
		debug:  debug,
		http:   newHTTPClient(insecure),
	}
}

// NewAnonymous creates a client for unauthenticated requests (e.g. login flow).
func NewAnonymous(server string, debug, insecure bool) *Client {
	return &Client{server: server, debug: debug, http: newHTTPClient(insecure)}
}

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

func (c *Client) maskToken() string {
	if len(c.token) <= 6 {
		return "***"
	}
	return c.token[:6] + "***"
}

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

type APIError struct {
	StatusCode int
}

func (e *APIError) Error() string {
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

func (c *Client) ListProviders() ([]string, error) {
	var providers []string
	if err := c.get("/api/v1/auth/providers", &providers); err != nil {
		return nil, err
	}
	return providers, nil
}

// PollToken checks whether the PAT for the given state is ready.
// Returns (nil, nil) when still pending (202), (token, nil) when ready (200).
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
		return nil, &APIError{StatusCode: resp.StatusCode}
	}
	var token models.UserToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (c *Client) GetProfile() (*models.User, error) {
	var u models.User
	if err := c.get("/api/v1/me/profile", &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *Client) ListUsers() ([]models.User, error) {
	var users []models.User
	if err := c.get("/api/v1/users", &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (c *Client) ListSessions(username string) ([]models.SSHSession, error) {
	var sessions []models.SSHSession
	if err := c.get(c.userPath(username)+"/sessions", &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}
