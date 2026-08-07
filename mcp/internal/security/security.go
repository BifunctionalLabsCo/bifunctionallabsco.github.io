package security

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/config"
)

const Confirmation = "CONFIRM"

var (
	sensitiveNames = regexp.MustCompile(`(?i)(^|/)(\.env($|\.)|id_[a-z0-9_-]+|credentials|secrets?|.*\.(pem|key|p12|pfx))$`)
	redactions     = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(api[_-]?key|token|secret|password)\s*[:=]\s*[^\s]+`),
		regexp.MustCompile(`gh[opsu]_[A-Za-z0-9_]{20,}`),
		regexp.MustCompile(`github_pat_[A-Za-z0-9_]{20,}`),
		regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----[\s\S]*?-----END [A-Z ]*PRIVATE KEY-----`),
	}
)

func Require(repo config.Repository, operation, confirmation string) error {
	if !repo.Allows(operation) {
		return fmt.Errorf("operation %q is not allowed for %s", operation, repo.Slug)
	}
	if strings.HasSuffix(operation, ":write") && confirmation != Confirmation {
		return fmt.Errorf("operation %q requires confirmation=%q", operation, Confirmation)
	}
	return nil
}

func ResolveExisting(root, relative string) (string, error) {
	rootReal, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", fmt.Errorf("resolve root: %w", err)
	}
	candidate := filepath.Join(rootReal, filepath.Clean(relative))
	real, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		return "", fmt.Errorf("resolve path: %w", err)
	}
	if err := ensureContained(rootReal, real); err != nil {
		return "", err
	}
	return real, nil
}

func ResolveForCreate(root, relative string) (string, error) {
	rootReal, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", fmt.Errorf("resolve root: %w", err)
	}
	candidate := filepath.Join(rootReal, filepath.Clean(relative))
	if err := ensureContained(rootReal, candidate); err != nil {
		return "", err
	}
	parent := filepath.Dir(candidate)
	for {
		if _, err := os.Lstat(parent); err == nil {
			realParent, err := filepath.EvalSymlinks(parent)
			if err != nil {
				return "", err
			}
			if err := ensureContained(rootReal, realParent); err != nil {
				return "", err
			}
			break
		}
		next := filepath.Dir(parent)
		if next == parent {
			return "", errors.New("could not resolve an existing parent")
		}
		parent = next
	}
	return candidate, nil
}

func ensureContained(root, candidate string) error {
	rel, err := filepath.Rel(root, candidate)
	if err != nil {
		return err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return errors.New("path escapes registered repository root")
	}
	return nil
}

func IsSensitivePath(path string) bool {
	return sensitiveNames.MatchString(filepath.ToSlash(path))
}

func Redact(content string) (string, int) {
	count := 0
	for _, pattern := range redactions {
		content = pattern.ReplaceAllStringFunc(content, func(string) string {
			count++
			return "[REDACTED]"
		})
	}
	return content, count
}
