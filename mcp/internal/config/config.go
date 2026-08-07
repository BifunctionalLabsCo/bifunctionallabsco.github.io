package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const (
	SensitivityPublic     = "public"
	SensitivityInternal   = "internal"
	SensitivityRestricted = "restricted"
)

type Config struct {
	Version          int          `json:"version"`
	Organization     string       `json:"organization"`
	AuditLog         string       `json:"auditLog"`
	MaxDocumentBytes int64        `json:"maxDocumentBytes"`
	CommandTimeout   string       `json:"commandTimeout"`
	BootstrapRoots   []string     `json:"bootstrapRoots"`
	Repositories     []Repository `json:"repositories"`
}

type Repository struct {
	Slug              string   `json:"slug"`
	LocalPath         string   `json:"localPath"`
	Visibility        string   `json:"visibility"`
	DefaultBranch     string   `json:"defaultBranch"`
	DocsRoots         []string `json:"docsRoots"`
	Sensitivity       string   `json:"sensitivity"`
	AllowedOperations []string `json:"allowedOperations"`
	ExcludedPaths     []string `json:"excludedPaths"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return nil, errors.New("decode config: trailing JSON content")
	}
	if err := cfg.Validate(filepath.Dir(path)); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) Validate(configDir string) error {
	if c.Version != 1 {
		return fmt.Errorf("unsupported config version %d", c.Version)
	}
	if strings.TrimSpace(c.Organization) == "" {
		return errors.New("organization is required")
	}
	if c.MaxDocumentBytes == 0 {
		c.MaxDocumentBytes = 1 << 20
	}
	if c.MaxDocumentBytes < 1024 || c.MaxDocumentBytes > 10<<20 {
		return errors.New("maxDocumentBytes must be between 1 KiB and 10 MiB")
	}
	if c.CommandTimeout == "" {
		c.CommandTimeout = "30s"
	}
	if _, err := time.ParseDuration(c.CommandTimeout); err != nil {
		return fmt.Errorf("invalid commandTimeout: %w", err)
	}
	if c.AuditLog == "" {
		c.AuditLog = filepath.Join(configDir, "var", "audit.jsonl")
	} else if !filepath.IsAbs(c.AuditLog) {
		c.AuditLog = filepath.Join(configDir, c.AuditLog)
	}
	for i, root := range c.BootstrapRoots {
		if !filepath.IsAbs(root) {
			root = filepath.Clean(filepath.Join(configDir, root))
		}
		info, err := os.Stat(root)
		if err != nil || !info.IsDir() {
			return fmt.Errorf("bootstrapRoot %q must be an existing directory", root)
		}
		resolved, err := filepath.EvalSymlinks(root)
		if err != nil {
			return fmt.Errorf("bootstrapRoot %q: %w", root, err)
		}
		c.BootstrapRoots[i] = resolved
	}

	seen := make(map[string]struct{}, len(c.Repositories))
	for i := range c.Repositories {
		r := &c.Repositories[i]
		if err := r.validate(c.Organization, configDir); err != nil {
			return fmt.Errorf("repository %d: %w", i, err)
		}
		if _, ok := seen[r.Slug]; ok {
			return fmt.Errorf("duplicate repository slug %q", r.Slug)
		}
		seen[r.Slug] = struct{}{}
	}
	return nil
}

func (r *Repository) validate(org, configDir string) error {
	parts := strings.Split(r.Slug, "/")
	if len(parts) != 2 || parts[0] != org || strings.TrimSpace(parts[1]) == "" {
		return fmt.Errorf("slug %q must be %s/name", r.Slug, org)
	}
	if !filepath.IsAbs(r.LocalPath) {
		r.LocalPath = filepath.Clean(filepath.Join(configDir, r.LocalPath))
	}
	resolved, err := filepath.EvalSymlinks(r.LocalPath)
	if err != nil {
		return fmt.Errorf("localPath %q: %w", r.LocalPath, err)
	}
	r.LocalPath = resolved
	info, err := os.Stat(r.LocalPath)
	if err != nil {
		return fmt.Errorf("localPath %q: %w", r.LocalPath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("localPath %q is not a directory", r.LocalPath)
	}
	if r.DefaultBranch == "" {
		r.DefaultBranch = "main"
	}
	if r.Sensitivity == "" {
		r.Sensitivity = SensitivityInternal
	}
	if !slices.Contains([]string{SensitivityPublic, SensitivityInternal, SensitivityRestricted}, r.Sensitivity) {
		return fmt.Errorf("invalid sensitivity %q", r.Sensitivity)
	}
	if len(r.DocsRoots) == 0 {
		return errors.New("at least one docsRoot is required")
	}
	for _, root := range append(slices.Clone(r.DocsRoots), r.ExcludedPaths...) {
		if filepath.IsAbs(root) || strings.HasPrefix(filepath.Clean(root), "..") {
			return fmt.Errorf("path %q must be relative and contained", root)
		}
	}
	return nil
}

func (r Repository) Allows(operation string) bool {
	return slices.Contains(r.AllowedOperations, operation)
}

func (c *Config) Repository(slug string) (Repository, bool) {
	for _, repo := range c.Repositories {
		if repo.Slug == slug {
			return repo, true
		}
	}
	return Repository{}, false
}

func (c *Config) Timeout() time.Duration {
	d, _ := time.ParseDuration(c.CommandTimeout)
	return d
}
