package server

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/audit"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/config"
	githubclient "github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/github"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/repos"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/scaffold"
)

func TestServerPublishesExpectedTools(t *testing.T) {
	root := t.TempDir()
	cfg := &config.Config{
		Organization:   "BifunctionalLabsCo",
		AuditLog:       filepath.Join(root, "audit.jsonl"),
		BootstrapRoots: []string{root},
		Repositories: []config.Repository{{
			Slug: "BifunctionalLabsCo/test", LocalPath: root, DocsRoots: []string{"docs"},
			AllowedOperations: []string{"repo:read", "docs:read", "issues:read", "issues:write", "pulls:read", "pulls:write", "branches:write"},
		}},
	}
	logger := audit.New(cfg.AuditLog)
	s := New(Dependencies{Config: cfg, Repos: repos.New(cfg), GitHub: githubclient.New(cfg, logger), Scaffold: scaffold.New(cfg, logger)})
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	ctx := context.Background()
	serverSession, err := s.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Close()
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clientSession.Close()

	want := map[string]bool{
		"list_repos": false, "search_docs": false, "read_doc": false,
		"search_issues": false, "create_issue": false, "update_issue": false,
		"run_socratic_intake": false, "generate_brand_direction": false,
		"generate_design_system": false, "bootstrap_repo": false, "publish_repo": false,
		"create_branch_scaffold": false, "generate_release_notes": false,
	}
	count := 0
	for tool, err := range clientSession.Tools(ctx, nil) {
		if err != nil {
			t.Fatal(err)
		}
		count++
		if _, ok := want[tool.Name]; ok {
			want[tool.Name] = true
		}
	}
	if count != 22 {
		t.Fatalf("tool count = %d, want 22", count)
	}
	for name, found := range want {
		if !found {
			t.Errorf("missing tool %s", name)
		}
	}
}
