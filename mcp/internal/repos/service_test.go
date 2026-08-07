package repos

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/config"
)

func TestSearchDocumentsUsesDocsRootsAndRedacts(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	doc := "---\ntitle: Platform security\nowner: bifunctional\nstatus: active\nupdated: 2026-08-07\ntags: [security]\n---\n\nThe token=supersecret is never returned.\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "security.md"), []byte(doc), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{MaxDocumentBytes: 1 << 20, Repositories: []config.Repository{{Slug: "BifunctionalLabsCo/test", LocalPath: root, DocsRoots: []string{"docs"}, AllowedOperations: []string{"docs:read"}}}}
	service := New(cfg)
	results, err := service.SearchDocuments("BifunctionalLabsCo/test", "platform", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Path != "docs/security.md" {
		t.Fatalf("unexpected results: %#v", results)
	}
	read, err := service.ReadDocument("BifunctionalLabsCo/test", "docs/security.md")
	if err != nil {
		t.Fatal(err)
	}
	if read.Redactions != 1 {
		t.Fatalf("redactions = %d, want 1", read.Redactions)
	}
	if _, err := service.ReadDocument("BifunctionalLabsCo/test", "README.md"); err == nil {
		t.Fatal("expected out-of-root read to fail")
	}
}
