# k8shell CLI

Command-line interface for managing k8shell resources — users, sessions, workspaces, and multi-environment contexts.

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
    username: alice          # populated automatically on login / context add
  - name: staging
    server: https://staging.k8shell.example.com
    token: <your-token>
    username: alice
    insecure: true           # skip TLS verification for this context
```

`username` and `tokenHash` are written automatically when a context is created via `login` or `context add`. Do not edit them by hand.

## Quick Start

```bash
# Log in via browser — saves a context automatically
k8shell login --server https://k8shell.example.com

# Log in again to the same server as a different user (--ignore bypasses the already-logged-in check)
k8shell login --server https://k8shell.example.com --ignore

# Or add a context manually with a PAT (token is prompted securely if --token is omitted)
k8shell context add prod --server https://k8shell.example.com

# Add a context for a server with a self-signed certificate
k8shell context add dev --server https://dev.k8shell.example.com --insecure

# List and switch contexts
k8shell context list
k8shell context use prod

# List users
k8shell user list

# List your active sessions
k8shell session list
```

## Documentation

Full documentation is available at [docs.k8shell.io/concepts/k8shell-cli](https://docs.k8shell.io/concepts/k8shell-cli).

## License

Copyright 2026 The k8shell CLI Authors. Licensed under the [GNU Affero General Public License v3.0](LICENSE).
