// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/common/pkg/userstr"
	"github.com/k8shell-io/k8shell/internal/client"
	"github.com/k8shell-io/k8shell/internal/config"
)

// Shared flags registered on workspace ssh / scp / code commands.
var (
	wsConnPod      string
	wsConnRepo     string
	wsConnRef      string
	wsConnHost     string
	wsConnB64      bool
	wsConnPrint    bool
	wsConnNoLookup bool
)

// wsConnResolve builds the SSH userstring, server hostname, and workspace
// owner username from the workspace identification flags. When b64Force is
// true the userstring is always base64-encoded regardless of --b64.
func wsConnResolve(b64Force bool) (ustr, server, username string, err error) {
	if wsConnPod == "" && wsConnRepo == "" {
		return "", "", "", fmt.Errorf("one of --pod or --repo is required")
	}
	if wsConnPod != "" && wsConnRepo != "" {
		return "", "", "", fmt.Errorf("--pod and --repo are mutually exclusive")
	}
	if wsConnRef != "" && wsConnPod != "" {
		return "", "", "", fmt.Errorf("--ref requires --repo")
	}

	ctx, err := cfg.ActiveContext()
	if err != nil {
		return "", "", "", err
	}

	var fields userstr.UserStrFields

	if wsConnPod != "" {
		if wsConnNoLookup {
			fields = userstr.UserStrFields{
				Username: ctx.Username,
				Pod:      wsConnPod,
			}
			username = ctx.Username
			server, err = wsServerFromContext(ctx)
			if err != nil {
				return "", "", "", err
			}
		} else {
			ws, err := client.New(ctx, debug, insecure || ctx.Insecure).GetWorkspace(wsConnPod)
			if err != nil {
				return "", "", "", err
			}
			fields = userstr.UserStrFields{
				Username:  ws.Username,
				Pod:       ws.Name,
				Namespace: ws.Namespace,
			}
			server = ws.ServerName
			username = ws.Username
		}
	} else {
		parts := strings.SplitN(wsConnRepo, "/", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return "", "", "", fmt.Errorf("--repo must be in owner/name format")
		}
		repoOwner, repoName := parts[0], parts[1]

		if wsConnNoLookup {
			username = ctx.Username
			if username == "" {
				username = repoOwner
			}
			fields = userstr.UserStrFields{
				Username:  username,
				RepoOwner: repoOwner,
				RepoName:  repoName,
				RepoRef:   wsConnRef,
			}
			server, err = wsServerFromContext(ctx)
			if err != nil {
				return "", "", "", err
			}
		} else {
			ws, err := wsFindByRepo(ctx, repoOwner, repoName, wsConnRef)
			if err != nil {
				return "", "", "", err
			}
			fields = userstr.UserStrFields{
				Username:  ws.Username,
				RepoOwner: repoOwner,
				RepoName:  repoName,
				RepoRef:   wsConnRef,
			}
			server = ws.ServerName
			username = ws.Username
		}
	}

	if wsConnHost != "" {
		server = wsConnHost
	}

	raw, err := fields.ToRawUserStr()
	if err != nil {
		return "", "", "", err
	}

	if b64Force || wsConnB64 {
		return "b64-" + base64.RawURLEncoding.EncodeToString([]byte(raw)), server, username, nil
	}
	return raw, server, username, nil
}

// wsFindByRepo lists workspaces and returns the first one matching the given
// repo owner, name, and optional ref.
func wsFindByRepo(ctx *config.Context, owner, name, ref string) (*models.WorkspaceDetails, error) {
	workspaces, err := client.New(ctx, debug, insecure || ctx.Insecure).ListWorkspaces("", false)
	if err != nil {
		return nil, err
	}
	for i := range workspaces {
		ws := &workspaces[i]
		if !strings.EqualFold(ws.RepoOwner, owner) || !strings.EqualFold(ws.RepoName, name) {
			continue
		}
		if ref != "" && ws.RepoRef != ref {
			continue
		}
		return ws, nil
	}
	if ref != "" {
		return nil, fmt.Errorf("no workspace found for repo %s/%s at ref %s", owner, name, ref)
	}
	return nil, fmt.Errorf("no workspace found for repo %s/%s", owner, name)
}

// wsServerFromContext derives the SSH hostname from the active context's server URL.
func wsServerFromContext(ctx *config.Context) (string, error) {
	u, err := url.Parse(ctx.Server)
	if err != nil || u.Hostname() == "" {
		return "", fmt.Errorf("cannot derive SSH host from context server %q; use --host", ctx.Server)
	}
	return u.Hostname(), nil
}
