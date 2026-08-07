package github

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/audit"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/config"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/security"
)

func TestMilestonesUsesGETAndStateOnlyIssueUpdate(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "args.log")
	installFakeGH(t, dir)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("GH_ARGS_LOG", logPath)
	cfg := &config.Config{Organization: "BifunctionalLabsCo", CommandTimeout: "5s", AuditLog: filepath.Join(dir, "audit.jsonl"), Repositories: []config.Repository{{Slug: "BifunctionalLabsCo/test", AllowedOperations: []string{"issues:read", "issues:write"}}}}
	client := New(cfg, audit.New(cfg.AuditLog))
	if _, err := client.Milestones(context.Background(), "BifunctionalLabsCo/test", "open"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.UpdateIssue(context.Background(), IssueUpdate{Repository: "BifunctionalLabsCo/test", Number: 7, State: "closed", Confirmation: security.Confirmation}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	log := string(b)
	if !strings.Contains(log, "api --method GET repos/BifunctionalLabsCo/test/milestones") {
		t.Errorf("milestone command was not GET: %s", log)
	}
	if !strings.Contains(log, "issue close 7 --repo BifunctionalLabsCo/test") {
		t.Errorf("state-only close was not executed: %s", log)
	}
}

func TestWriteRequiresConfirmationBeforeCommand(t *testing.T) {
	dir := t.TempDir()
	installFakeGH(t, dir)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	cfg := &config.Config{CommandTimeout: "5s", AuditLog: filepath.Join(dir, "audit.jsonl"), Repositories: []config.Repository{{Slug: "BifunctionalLabsCo/test", AllowedOperations: []string{"issues:write"}}}}
	client := New(cfg, audit.New(cfg.AuditLog))
	if _, err := client.CreateIssue(context.Background(), IssueCreate{Repository: "BifunctionalLabsCo/test", Title: "No confirmation"}); err == nil {
		t.Fatal("expected confirmation failure")
	}
}

func installFakeGH(t *testing.T, dir string) {
	t.Helper()
	path := filepath.Join(dir, "gh")
	script := `#!/bin/sh
printf '%s\n' "$*" >> "$GH_ARGS_LOG"
case "$1 $2" in
  "api --method") printf '[]\n' ;;
  "issue view") printf '{"number":7,"state":"CLOSED"}\n' ;;
  *) printf 'https://github.test/result\n' ;;
esac
`
	if err := os.WriteFile(path, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
}
