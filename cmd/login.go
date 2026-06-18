package cmd

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/client"
	"github.com/k8shell-io/k8shell/internal/config"
	"github.com/spf13/cobra"
)

var (
	loginServer  string
	loginName    string
	loginTimeout time.Duration
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in via browser and save credentials to a context",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, ctx := range cfg.Contexts {
			if ctx.Server == loginServer && ctx.Token != "" {
				c := client.New(&ctx, debug, insecure || ctx.Insecure)
				if profile, err := c.GetProfile(); err == nil {
					fmt.Printf("Already logged in as %s (context %q).\n", profile.Username, ctx.Name)
					return nil
				}
				break
			}
		}

		state, err := randomState()
		if err != nil {
			return fmt.Errorf("generating state: %w", err)
		}

		c := client.NewAnonymous(loginServer, debug, insecure)

		providers, err := c.ListProviders()
		if err != nil {
			return fmt.Errorf("fetching providers: %w", err)
		}
		if len(providers) == 0 {
			return fmt.Errorf("no login providers configured on server")
		}

		provider, err := pickProvider(providers)
		if err != nil {
			return err
		}

		loginURL := buildLoginURL(loginServer, state, provider)
		fmt.Println("Opening browser for login...")
		if err := openBrowser(loginURL); err != nil {
			fmt.Printf("Could not open browser automatically. Open this URL:\n\n  %s\n\n", loginURL)
		}

		fmt.Print("Waiting for login to complete...")
		token, err := pollForToken(c, state, loginTimeout)
		if err != nil {
			fmt.Println()
			return err
		}
		fmt.Println(" done.")

		name := loginName
		if name == "" {
			if u, err := url.Parse(loginServer); err == nil {
				name = u.Hostname()
			} else {
				name = loginServer
			}
		}

		_ = cfg.DeleteContext(name)
		if err := cfg.AddContext(config.Context{Name: name, Server: loginServer, Token: token.Token, Insecure: insecure}); err != nil {
			return err
		}
		cfg.CurrentContext = name
		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Printf("Logged in as %s. Context %q saved and set as active.\n", token.Username, name)
		return nil
	},
}

func init() {
	loginCmd.Flags().StringVar(&loginServer, "server", "", "API server URL (required)")
	loginCmd.Flags().StringVar(&loginName, "name", "", "context name (defaults to server hostname)")
	loginCmd.Flags().DurationVar(&loginTimeout, "timeout", 5*time.Minute, "time to wait for browser login")
	_ = loginCmd.MarkFlagRequired("server")
}

func pollForToken(c *client.Client, state string, timeout time.Duration) (*models.UserToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("login timed out after %s", timeout)
		case <-ticker.C:
			token, err := c.PollToken(state)
			if err != nil {
				return nil, err
			}
			if token != nil {
				return token, nil
			}
		}
	}
}

func buildLoginURL(server, state, provider string) string {
	redirectURI, _ := url.Parse(server + "/auth/callback")
	rq := redirectURI.Query()
	rq.Set("provider", provider)
	rq.Set("redirect", "/")
	redirectURI.RawQuery = rq.Encode()

	u, _ := url.Parse(server + "/api/v1/auth/login")
	q := u.Query()
	q.Set("createPAT", "true")
	q.Set("state", state)
	q.Set("provider", provider)
	q.Set("redirect_uri", redirectURI.String())
	u.RawQuery = q.Encode()
	return u.String()
}

func pickProvider(providers []string) (string, error) {
	if len(providers) == 1 {
		return providers[0], nil
	}
	fmt.Println("Available login providers:")
	for i, p := range providers {
		// strip common prefix for readability (e.g. "idp.k8shell.io/github" → "github")
		label := p
		if idx := strings.LastIndex(p, "/"); idx >= 0 {
			label = p[idx+1:]
		}
		fmt.Printf("  %d) %s\n", i+1, label)
	}
	fmt.Print("Select provider [1]: ")
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return providers[0], nil
	}
	n, err := strconv.Atoi(line)
	if err != nil || n < 1 || n > len(providers) {
		return "", fmt.Errorf("invalid selection %q", line)
	}
	return providers[n-1], nil
}

func randomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func openBrowser(rawURL string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", rawURL)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", rawURL)
	default:
		cmd = exec.Command("xdg-open", rawURL)
	}
	return cmd.Start()
}
