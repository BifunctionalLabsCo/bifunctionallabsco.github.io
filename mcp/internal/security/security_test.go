package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/config"
)

func TestResolveExistingRejectsTraversalAndSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret.md"), []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ResolveExisting(root, filepath.Join("..", filepath.Base(outside), "secret.md")); err == nil {
		t.Fatal("expected traversal to fail")
	}
	if err := os.Symlink(outside, filepath.Join(root, "escape")); err != nil {
		t.Fatal(err)
	}
	if _, err := ResolveExisting(root, "escape/secret.md"); err == nil {
		t.Fatal("expected symlink escape to fail")
	}
}

func TestRedact(t *testing.T) {
	input := "api_key=very-secret\ntoken: abcdef\n" + "github_" + "pat_abcdefghijklmnopqrstuvwxyz123456"
	out, count := Redact(input)
	if count != 3 {
		t.Fatalf("redactions = %d, want 3", count)
	}
	if strings.Contains(out, "very-secret") || strings.Contains(out, "github_pat_") {
		t.Fatalf("sensitive content survived: %s", out)
	}
}

func TestRequireWriteConfirmation(t *testing.T) {
	repo := config.Repository{Slug: "BifunctionalLabsCo/test", AllowedOperations: []string{"issues:write"}}
	if err := Require(repo, "issues:write", ""); err == nil {
		t.Fatal("expected missing confirmation to fail")
	}
	if err := Require(repo, "issues:write", Confirmation); err != nil {
		t.Fatalf("confirmed write failed: %v", err)
	}
}
