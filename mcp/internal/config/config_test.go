package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateRejectsDuplicateAndEscapingPaths(t *testing.T) {
	root := t.TempDir()
	cfg := Config{
		Version:      1,
		Organization: "BifunctionalLabsCo",
		Repositories: []Repository{{Slug: "BifunctionalLabsCo/test", LocalPath: root, DocsRoots: []string{"../private"}}},
	}
	if err := cfg.Validate(root); err == nil {
		t.Fatal("expected escaping docs root to fail")
	}
}

func TestLoadDisallowsUnknownFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "repos.json")
	if err := os.WriteFile(path, []byte(`{"version":1,"organization":"BifunctionalLabsCo","unknown":true,"repositories":[]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected unknown field to fail")
	}
}
