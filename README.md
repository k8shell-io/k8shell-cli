# k8shell CLI

Command-line interface for managing k8shell resources — users, sessions, and multi-environment contexts.

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

### login

```
k8shell login --server <url> [--name <context-name>] [--timeout <duration>]
```

Authenticates via browser OAuth and saves the resulting PAT to a named context. The context name defaults to the server hostname. Pass `--insecure` (global flag) to skip TLS verification; this is persisted in the saved context.

### context

Aliases: `ctx`

| Command | Description |
|---|---|
| `context list [--sort <fields>]` | List all configured contexts |
| `context add <name> --server <url> [--insecure]` | Add a new context; token prompted if `--token` is omitted |
| `context use <name>` | Switch the active context |
| `context delete <name>` | Remove a context |

`--insecure` on `context add` is stored in the config and applied automatically to every command that uses that context — you won't need to pass `--insecure` again.

### user

Aliases: `usr`

| Command | Description |
|---|---|
| `user list [--sort <fields>]` | List all users |
| `user session list [-u <username>] [--sort <fields>]` | List sessions; defaults to the authenticated user |

`list` accepts the alias `ls`. Sort fields are the column names shown in the output, prefixed with `-` for descending order (e.g. `--sort username,-email`).

## Global Flags

| Flag | Description |
|---|---|
| `--config <path>` | Path to config file |
| `-c, --context <name>` | Override the active context for this invocation |
| `--insecure` | Skip TLS certificate verification |
| `--json` | Output as JSON |
| `--no-ansi` | Disable color output |
| `-w, --wrap` | Allow lines to wrap beyond terminal width |
| `--debug` | Print request/response headers to stderr |

## License

Copyright 2026 The k8shell CLI Authors. Licensed under the [GNU Affero General Public License v3.0](LICENSE).
