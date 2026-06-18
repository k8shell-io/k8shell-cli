# k8shell CLI

Command-line interface for managing k8shell resources — users, sessions, and multi-environment contexts.

> Requires a running k8shell platform with the API server reachable at a configured address.

## Installation

```bash
go install github.com/k8shell-io/k8shell@latest
```

Or build from source:

```bash
git clone https://github.com/k8shell-io/k8shell-cli.git
cd k8shell-cli
go build -o k8shell .
```

## Configuration

The config file lives at `~/.config/k8shell/config.yaml` by default. Override with `--config <path>`.

```yaml
current-context: production
contexts:
  - name: production
    server: https://k8shell.example.com
    token: <your-token>
  - name: staging
    server: https://staging.k8shell.example.com
    token: <your-token>
    insecure: true   # skip TLS verification for this context
```

## Quick Start

```bash
# Log in via browser — saves a context automatically
k8shell login --server https://k8shell.example.com

# Or add a context manually (token is prompted securely if --token is omitted)
k8shell context add prod --server https://k8shell.example.com

# Add a context for a server with a self-signed certificate
k8shell context add dev --server https://dev.k8shell.example.com --insecure

# List all configured contexts
k8shell context list

# Switch active context
k8shell context use prod

# List users
k8shell user list

# List your active sessions
k8shell user session list

# List sessions for a specific user
k8shell user session list --user alice
```

## Commands

For a complete reference of all commands, see the [k8shell CLI docs](https://docs.k8shell.io/concepts/k8shell-cli).

## License

Copyright 2026 The k8shell CLI Authors. Licensed under the [GNU Affero General Public License v3.0](LICENSE).
