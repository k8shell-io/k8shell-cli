# k8shell-cli

CLI tool for managing k8shell resources. Written in Go using [cobra](https://github.com/spf13/cobra).

## Project structure

```
main.go                  # entrypoint — calls cmd.Execute()
cmd/                     # cobra commands
  root.go                # root command, global flags, PersistentPreRunE
  context.go             # `context` subcommand group
  context_add.go
  context_delete.go
  context_list.go
  context_use.go
  user.go                # `user` subcommand group
  user_list.go
  user_session.go        # `user session` subcommand group
  user_session_list.go
internal/
  client/client.go       # HTTP client — wraps API calls (ListUsers, ListSessions)
  config/config.go       # YAML config: load, save, context CRUD
  output/output.go       # table/JSON printer with ANSI color support
```

## Config

Default path: `~/.config/k8shell/config.yaml`

```yaml
current-context: prod
contexts:
  - name: prod
    server: https://k8shell.example.com
    token: <pat>
```

`config.Config` holds the context list. `ActiveContext()` resolves the current one. Config is saved back to disk after any mutation.

## Key conventions

- All commands return errors — no `os.Exit` in command handlers, only in `Execute()`.
- Output goes through `output.Printer` (table or JSON). Never write directly to stdout in command handlers — use `printer`.
- The HTTP client is constructed per-command from the active context: `client.New(ctx, debug)`.
- Global flags (`--json`, `--no-ansi`, `--wrap`, `--debug`, `--context`, `--config`) are defined in `root.go` and wired in `PersistentPreRunE`.
- New top-level command groups follow the pattern in `cmd/context.go` + `cmd/user.go`: define a parent `cobra.Command`, add subcommands in `init()`, register with `rootCmd` in `root.go`'s `init()`.

## Build & run

```bash
go build -o k8shell .
./k8shell --help
```

## Module

`github.com/k8shell-io/k8shell` — shared models come from `github.com/k8shell-io/common`.
