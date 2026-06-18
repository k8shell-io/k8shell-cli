// Copyright 2026 The k8shell CLI Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package client

import "github.com/k8shell-io/common/pkg/models"

// GetProfile returns the profile of the authenticated user.
func (c *Client) GetProfile() (*models.User, error) {
	var u models.User
	if err := c.get("/api/v1/me/profile", &u); err != nil {
		return nil, err
	}
	return &u, nil
}

// ListUsers returns all users visible to the authenticated token.
func (c *Client) ListUsers() ([]models.User, error) {
	var users []models.User
	if err := c.get("/api/v1/users", &users); err != nil {
		return nil, err
	}
	return users, nil
}

// ListSessions returns SSH sessions for the given username, or for the authenticated user if username is empty.
func (c *Client) ListSessions(username string) ([]models.SSHSession, error) {
	var sessions []models.SSHSession
	if err := c.get(c.userPath(username)+"/sessions", &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}
