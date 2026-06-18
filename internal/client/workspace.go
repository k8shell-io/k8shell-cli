// Copyright 2026 The k8shell CLI Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package client

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/k8shell-io/common/pkg/models"
)

// WorkspaceCreateRequest is the payload for creating a workspace.
// Either Blueprint or (RepoOwner + RepoName) must be set.
type WorkspaceCreateRequest struct {
	Username  string `json:"username"`
	Blueprint string `json:"blueprint,omitempty"`
	RepoOwner string `json:"repoOwner,omitempty"`
	RepoName  string `json:"repoName,omitempty"`
	RepoRef   string `json:"repoRef,omitempty"`
}

// WorkspaceCreateResponse is the 202 Accepted body returned by POST /workspaces.
type WorkspaceCreateResponse struct {
	Workspace  string `json:"workspace"`
	JobID      string `json:"jobId"`
	MonitorURL string `json:"monitorUrl"`
}

// ListWorkspaces returns workspaces visible to the authenticated token.
// When username is non-empty it is passed as a query parameter to filter by owner.
func (c *Client) ListWorkspaces(username string) ([]models.WorkspaceDetails, error) {
	path := "/api/v1/workspaces"
	if username != "" {
		path += "?username=" + url.QueryEscape(username)
	}
	var workspaces []models.WorkspaceDetails
	if err := c.get(path, &workspaces); err != nil {
		return nil, err
	}
	return workspaces, nil
}

// CreateWorkspace submits a workspace creation request and returns the 202 response.
func (c *Client) CreateWorkspace(req WorkspaceCreateRequest) (*WorkspaceCreateResponse, error) {
	var resp WorkspaceCreateResponse
	if err := c.post("/api/v1/workspaces", req, &resp); err != nil {
		return nil, err
	}
	// The server may return a monitorUrl with the jobId duplicated (e.g. /jobs/{id}/{id}).
	// Strip the trailing duplicate so the SSE request goes to the correct path.
	if resp.JobID != "" && strings.Count(resp.MonitorURL, resp.JobID) > 1 {
		if idx := strings.LastIndex(resp.MonitorURL, "/"+resp.JobID); idx >= 0 {
			resp.MonitorURL = resp.MonitorURL[:idx]
		}
	}
	return &resp, nil
}

// MonitorWorkspace opens an SSE stream at monitorURL and returns the response body for the caller to read.
// monitorURL may be a full URL or a path; a path is resolved against the client's server base.
func (c *Client) MonitorWorkspace(monitorURL string) (io.ReadCloser, error) {
	if u, err := url.Parse(monitorURL); err == nil && u.Scheme == "" {
		monitorURL = c.server + monitorURL
	}
	req, err := http.NewRequest(http.MethodGet, monitorURL, nil)
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "text/event-stream")

	if c.debug {
		c.debugRequest(req)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if c.debug {
		c.debugResponse(resp)
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, &APIError{StatusCode: resp.StatusCode}
	}
	return resp.Body, nil
}

// GetWorkspace returns the details of the named workspace.
func (c *Client) GetWorkspace(name string) (*models.WorkspaceDetails, error) {
	var ws models.WorkspaceDetails
	if err := c.get("/api/v1/workspaces/"+name, &ws); err != nil {
		return nil, err
	}
	return &ws, nil
}

// DeleteWorkspace shuts down the named workspace.
// When deleteData is true, ?delete=true is appended to request permanent data deletion.
func (c *Client) DeleteWorkspace(name string, deleteData bool) error {
	path := "/api/v1/workspaces/" + name
	if deleteData {
		path += "?delete=true"
	}
	return c.delete(path)
}
