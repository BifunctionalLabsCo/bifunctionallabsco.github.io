package scaffold

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/audit"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/config"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/ideation"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/security"
)

func TestGenerateCreatesCompleteRepo(t *testing.T) {
	root := t.TempDir()
	// The scaffold intentionally uses the operator's Git identity. Supply one
	// only for this isolated test so CI does not depend on runner configuration.
	t.Setenv("GIT_AUTHOR_NAME", "Bifunctional MCP Test")
	t.Setenv("GIT_AUTHOR_EMAIL", "mcp-test@bifunctional.invalid")
	t.Setenv("GIT_COMMITTER_NAME", "Bifunctional MCP Test")
	t.Setenv("GIT_COMMITTER_EMAIL", "mcp-test@bifunctional.invalid")
	cfg := &config.Config{Organization: "BifunctionalLabsCo", BootstrapRoots: []string{root}, AuditLog: filepath.Join(root, "audit.jsonl")}
	service := New(cfg, audit.New(cfg.AuditLog))
	result, err := service.Generate(Request{Project: completeProject(), Parent: root, License: "proprietary", Gitignore: "go", Confirmation: security.Confirmation})
	if err != nil {
		t.Fatal(err)
	}
	for _, relative := range []string{"README.md", "AGENTS.md", "CONTRIBUTING.md", "LICENSE", "docs/index.md", "docs/design/brand-direction.md", "docs/design/design-system.md", ".git"} {
		if _, err := os.Stat(filepath.Join(result.Path, relative)); err != nil {
			t.Errorf("missing %s: %v", relative, err)
		}
	}
	if err := exec.Command("git", "-C", result.Path, "rev-parse", "--verify", "HEAD").Run(); err != nil {
		t.Fatalf("expected initial commit: %v", err)
	}
	if _, err := service.Generate(Request{Project: completeProject(), Parent: root, Confirmation: security.Confirmation}); err == nil {
		t.Fatal("expected existing destination to fail")
	}
}

func TestGenerateRequiresConfirmationAndAllowedRoot(t *testing.T) {
	root := t.TempDir()
	service := New(&config.Config{BootstrapRoots: []string{root}}, audit.New(filepath.Join(root, "audit.jsonl")))
	if _, err := service.Generate(Request{Project: completeProject(), Parent: root}); err == nil {
		t.Fatal("expected confirmation failure")
	}
	if _, err := service.Generate(Request{Project: completeProject(), Parent: t.TempDir(), Confirmation: security.Confirmation}); err == nil {
		t.Fatal("expected root policy failure")
	}
}

func completeProject() ideation.Project {
	return ideation.Project{Name: "Test Project", Purpose: "Help operators test safely.", Audience: "Platform operators", ProductType: "CLI", Success: "A verified scaffold", InheritedTraits: []string{"precise"}, DistinctTraits: []string{"calm"}, InitialWorkflows: []string{"Create project"}, Constraints: []string{"Local only"}}
}
