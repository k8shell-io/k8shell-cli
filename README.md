# k8shell

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
```

## Quick Start

```bash
# Add your first context (token is prompted securely)
k8shell context add prod --server https://k8shell.example.com

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

### context

| Command | Description |
|---|---|
| `context list` | List all configured contexts |
| `context add <name> --server <url>` | Add a new context (token prompted if omitted) |
| `context use <name>` | Switch the active context |
| `context delete <name>` | Remove a context |

### user

| Command | Description |
|---|---|
| `user list` | List all users |
| `user session list` | List sessions (defaults to authenticated user) |
| `user session list --user <name>` | List sessions for a specific user |

## Global Flags

| Flag | Description |
|---|---|
| `--config <path>` | Path to config file |
| `-c, --context <name>` | Override the active context for this invocation |
| `--json` | Output as JSON |
| `--no-ansi` | Disable color output |
| `-w, --wrap` | Allow lines to wrap beyond terminal width |
| `--debug` | Print request/response headers to stderr |

## License

See [LICENSE](LICENSE).
