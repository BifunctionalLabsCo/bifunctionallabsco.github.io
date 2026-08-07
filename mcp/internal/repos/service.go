package repos

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/config"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/okf"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/security"
)

type Service struct {
	cfg *config.Config
}

type Summary struct {
	Slug              string   `json:"slug"`
	Visibility        string   `json:"visibility"`
	DefaultBranch     string   `json:"defaultBranch"`
	DocsRoots         []string `json:"docsRoots"`
	Sensitivity       string   `json:"sensitivity"`
	AllowedOperations []string `json:"allowedOperations"`
	Available         bool     `json:"available"`
}

type Metadata struct {
	Slug          string `json:"slug"`
	LocalPath     string `json:"localPath"`
	CurrentBranch string `json:"currentBranch"`
	Head          string `json:"head"`
	Dirty         bool   `json:"dirty"`
	Remote        string `json:"remote,omitempty"`
}

type SearchResult struct {
	Repository string       `json:"repository"`
	Path       string       `json:"path"`
	Title      string       `json:"title,omitempty"`
	Score      int          `json:"score"`
	Snippet    string       `json:"snippet"`
	Metadata   okf.Metadata `json:"metadata"`
}

func New(cfg *config.Config) *Service { return &Service{cfg: cfg} }

func (s *Service) List() []Summary {
	result := make([]Summary, 0, len(s.cfg.Repositories))
	for _, repo := range s.cfg.Repositories {
		_, err := os.Stat(repo.LocalPath)
		result = append(result, Summary{
			Slug: repo.Slug, Visibility: repo.Visibility, DefaultBranch: repo.DefaultBranch,
			DocsRoots: repo.DocsRoots, Sensitivity: repo.Sensitivity,
			AllowedOperations: repo.AllowedOperations, Available: err == nil,
		})
	}
	return result
}

func (s *Service) Get(slug string) (config.Repository, error) {
	repo, ok := s.cfg.Repository(slug)
	if !ok {
		return config.Repository{}, fmt.Errorf("repository %q is not registered", slug)
	}
	return repo, nil
}

func (s *Service) Metadata(ctx context.Context, slug string) (Metadata, error) {
	repo, err := s.Get(slug)
	if err != nil {
		return Metadata{}, err
	}
	if err := security.Require(repo, "repo:read", ""); err != nil {
		return Metadata{}, err
	}
	run := func(args ...string) (string, error) {
		cmd := exec.CommandContext(ctx, "git", append([]string{"-C", repo.LocalPath}, args...)...)
		out, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
		}
		return strings.TrimSpace(string(out)), nil
	}
	branch, err := run("branch", "--show-current")
	if err != nil {
		return Metadata{}, err
	}
	head, err := run("rev-parse", "HEAD")
	if err != nil {
		return Metadata{}, err
	}
	status, err := run("status", "--porcelain")
	if err != nil {
		return Metadata{}, err
	}
	remote, _ := run("remote", "get-url", "origin")
	return Metadata{Slug: slug, LocalPath: repo.LocalPath, CurrentBranch: branch, Head: head, Dirty: status != "", Remote: remote}, nil
}

func (s *Service) ReadDocument(slug, relative string) (okf.Document, error) {
	repo, err := s.Get(slug)
	if err != nil {
		return okf.Document{}, err
	}
	if err := security.Require(repo, "docs:read", ""); err != nil {
		return okf.Document{}, err
	}
	if !isAllowedDocument(repo, relative) {
		return okf.Document{}, errors.New("document is outside configured docs roots or excluded")
	}
	path, err := security.ResolveExisting(repo.LocalPath, relative)
	if err != nil {
		return okf.Document{}, err
	}
	if security.IsSensitivePath(path) {
		return okf.Document{}, errors.New("sensitive file type is not readable")
	}
	info, err := os.Stat(path)
	if err != nil {
		return okf.Document{}, err
	}
	if info.Size() > s.cfg.MaxDocumentBytes {
		return okf.Document{}, fmt.Errorf("document exceeds %d byte limit", s.cfg.MaxDocumentBytes)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return okf.Document{}, err
	}
	redacted, count := security.Redact(string(b))
	meta, body, err := okf.Parse(redacted)
	if err != nil {
		return okf.Document{}, err
	}
	return okf.Document{Path: filepath.ToSlash(relative), Metadata: meta, Content: body, Redactions: count}, nil
}

func (s *Service) SearchDocuments(slug, query string, limit int) ([]SearchResult, error) {
	if strings.TrimSpace(query) == "" {
		return nil, errors.New("query is required")
	}
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	repo, err := s.Get(slug)
	if err != nil {
		return nil, err
	}
	if err := security.Require(repo, "docs:read", ""); err != nil {
		return nil, err
	}
	repoRoot, err := security.ResolveExisting(repo.LocalPath, ".")
	if err != nil {
		return nil, err
	}
	terms := strings.Fields(strings.ToLower(query))
	var results []SearchResult
	for _, docsRoot := range repo.DocsRoots {
		root, err := security.ResolveExisting(repo.LocalPath, docsRoot)
		if err != nil {
			continue
		}
		err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return nil
			}
			rel, err := filepath.Rel(repoRoot, path)
			if err != nil || excluded(repo, rel) {
				if entry.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if entry.IsDir() || !isMarkdown(path) || security.IsSensitivePath(path) {
				return nil
			}
			doc, err := s.ReadDocument(slug, rel)
			if err != nil {
				return nil
			}
			haystack := strings.ToLower(doc.Metadata.Title + "\n" + strings.Join(doc.Metadata.Tags, " ") + "\n" + doc.Content)
			score := 0
			for _, term := range terms {
				score += strings.Count(haystack, term)
				if strings.Contains(strings.ToLower(doc.Metadata.Title), term) {
					score += 5
				}
			}
			if score > 0 {
				results = append(results, SearchResult{Repository: slug, Path: filepath.ToSlash(rel), Title: doc.Metadata.Title, Score: score, Snippet: snippet(doc.Content, terms), Metadata: doc.Metadata})
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Path < results[j].Path
		}
		return results[i].Score > results[j].Score
	})
	if len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

func isAllowedDocument(repo config.Repository, relative string) bool {
	clean := filepath.Clean(relative)
	if excluded(repo, clean) || !isMarkdown(clean) {
		return false
	}
	for _, root := range repo.DocsRoots {
		rel, err := filepath.Rel(filepath.Clean(root), clean)
		if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

func excluded(repo config.Repository, relative string) bool {
	clean := filepath.ToSlash(filepath.Clean(relative))
	for _, excluded := range repo.ExcludedPaths {
		excluded = strings.TrimSuffix(filepath.ToSlash(filepath.Clean(excluded)), "/")
		if clean == excluded || strings.HasPrefix(clean, excluded+"/") {
			return true
		}
	}
	return false
}

func isMarkdown(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".mdx"
}

func snippet(content string, terms []string) string {
	compact := strings.Join(strings.Fields(content), " ")
	lower := strings.ToLower(compact)
	index := 0
	for _, term := range terms {
		if i := strings.Index(lower, term); i >= 0 {
			index = i
			break
		}
	}
	start := max(0, index-80)
	end := min(len(compact), index+240)
	if start > 0 {
		compact = "..." + compact[start:end]
	} else {
		compact = compact[start:end]
	}
	if end < len(lower) {
		compact += "..."
	}
	return compact
}
