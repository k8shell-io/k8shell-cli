// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package config

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Context holds the connection details for a single k8shell server.
type Context struct {
	Name      string `yaml:"name"`
	Server    string `yaml:"server"`
	Token     string `yaml:"token"`
	Username  string `yaml:"username,omitempty"`
	TokenHash string `yaml:"tokenHash,omitempty"`
	Insecure  bool   `yaml:"insecure,omitempty"`
}

// SetIntegrity computes and stores a hash of the token and username.
// Call this after setting Token and Username before saving the context.
func (ctx *Context) SetIntegrity() {
	ctx.TokenHash = tokenHash(ctx.Token, ctx.Username)
}

// tokenHash returns hex(SHA256(token + "\x00" + username)).
func tokenHash(token, username string) string {
	h := sha256.Sum256([]byte(token + "\x00" + username))
	return hex.EncodeToString(h[:])
}

// verifyIntegrity checks the stored hash against the current token and username.
// If no hash is stored the check is skipped; a present but wrong hash is always an error.
func (ctx *Context) verifyIntegrity() error {
	if ctx.TokenHash == "" {
		return nil
	}
	if tokenHash(ctx.Token, ctx.Username) != ctx.TokenHash {
		return fmt.Errorf("context %q integrity check failed: token/username mismatch — re-run 'k8shell login' to refresh", ctx.Name)
	}
	return nil
}

// Config is the top-level configuration loaded from the YAML config file.
type Config struct {
	CurrentContext string    `yaml:"current-context"`
	Contexts       []Context `yaml:"contexts"`
	path           string
}

// ActiveContext returns the context matching CurrentContext, or an error if none is set or found.
func (c *Config) ActiveContext() (*Context, error) {
	if c.CurrentContext == "" {
		return nil, fmt.Errorf("no active context; set current-context in config or use --context")
	}
	for i := range c.Contexts {
		if c.Contexts[i].Name == c.CurrentContext {
			if err := c.Contexts[i].verifyIntegrity(); err != nil {
				return nil, err
			}
			return &c.Contexts[i], nil
		}
	}
	return nil, fmt.Errorf("context %q not found in config", c.CurrentContext)
}

// ActiveContextByName returns the context with the given name, or an error if not found.
func (c *Config) ActiveContextByName(name string) (*Context, error) {
	for i := range c.Contexts {
		if c.Contexts[i].Name == name {
			return &c.Contexts[i], nil
		}
	}
	return nil, fmt.Errorf("context %q not found", name)
}

// AddContext appends ctx to the config, returning an error if a context with the same name already exists.
func (c *Config) AddContext(ctx Context) error {
	for _, existing := range c.Contexts {
		if existing.Name == ctx.Name {
			return fmt.Errorf("context %q already exists", ctx.Name)
		}
	}
	c.Contexts = append(c.Contexts, ctx)
	return nil
}

// DeleteContext removes the named context from the config and clears CurrentContext if it matched.
func (c *Config) DeleteContext(name string) error {
	for i, ctx := range c.Contexts {
		if ctx.Name == name {
			c.Contexts = append(c.Contexts[:i], c.Contexts[i+1:]...)
			if c.CurrentContext == name {
				c.CurrentContext = ""
			}
			return nil
		}
	}
	return fmt.Errorf("context %q not found", name)
}

// Save serialises the config to its YAML file, creating parent directories as needed.
func (c *Config) Save() error {
	if err := os.MkdirAll(filepath.Dir(c.path), 0700); err != nil {
		return err
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0600)
}

// Load reads the config file at path (or the default path if empty) and returns a parsed Config.
// A missing file is not an error; an empty Config is returned instead.
func Load(path string) (*Config, error) {
	if path == "" {
		var err error
		path, err = defaultPath()
		if err != nil {
			return nil, err
		}
	}

	cfg := &Config{path: path}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	return cfg, nil
}

// defaultPath returns ~/.config/k8shell/config.yaml.
func defaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "k8shell", "config.yaml"), nil
}
