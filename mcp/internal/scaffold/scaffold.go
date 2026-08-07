package scaffold

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/audit"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/config"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/ideation"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/security"
)

var safeName = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{1,62}$`)

type Service struct {
	cfg   *config.Config
	audit *audit.Logger
}

type Request struct {
	Project      ideation.Project `json:"project"`
	Parent       string           `json:"parent"`
	License      string           `json:"license"`
	Gitignore    string           `json:"gitignore"`
	Confirmation string           `json:"confirmation"`
}

type Result struct {
	Path  string   `json:"path"`
	Files []string `json:"files"`
}

func New(cfg *config.Config, logger *audit.Logger) *Service {
	return &Service{cfg: cfg, audit: logger}
}

func Preview(req Request) (Result, error) {
	if err := ideation.Validate(req.Project); err != nil {
		return Result{}, err
	}
	if err := validateOptions(req); err != nil {
		return Result{}, err
	}
	brief, _ := ideation.ProjectBrief(req.Project)
	brand, _ := ideation.BrandDirection(req.Project)
	design, _ := ideation.DesignSystem(req.Project)
	files := projectFiles(req, brief, brand, design)
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return Result{Path: filepath.Join(req.Parent, slugify(req.Project.Name)), Files: paths}, nil
}

func (s *Service) Generate(req Request) (Result, error) {
	if req.Confirmation != security.Confirmation {
		return Result{}, fmt.Errorf("scaffolding requires confirmation=%q", security.Confirmation)
	}
	if err := ideation.Validate(req.Project); err != nil {
		return Result{}, err
	}
	if err := validateOptions(req); err != nil {
		return Result{}, err
	}
	name := slugify(req.Project.Name)
	if !safeName.MatchString(name) {
		return Result{}, errors.New("project name cannot produce a safe repository slug")
	}
	parent, err := s.allowedParent(req.Parent)
	if err != nil {
		return Result{}, err
	}
	destination, err := security.ResolveForCreate(parent, name)
	if err != nil {
		return Result{}, err
	}
	if _, err := os.Lstat(destination); err == nil {
		return Result{}, ideation.ErrExists
	} else if !os.IsNotExist(err) {
		return Result{}, err
	}
	if err := s.audit.Record(audit.Event{Operation: "repo:scaffold", Target: destination, Outcome: "requested"}); err != nil {
		return Result{}, fmt.Errorf("write blocked because audit log is unavailable: %w", err)
	}
	brief, _ := ideation.ProjectBrief(req.Project)
	brand, _ := ideation.BrandDirection(req.Project)
	design, _ := ideation.DesignSystem(req.Project)
	files := projectFiles(req, brief, brand, design)
	if err := writeAtomically(destination, files); err != nil {
		return Result{}, err
	}
	if err := run(destination, "git", "init", "-b", "main"); err != nil {
		_ = os.RemoveAll(destination)
		return Result{}, err
	}
	if err := run(destination, "git", "add", "--all"); err != nil {
		_ = os.RemoveAll(destination)
		return Result{}, err
	}
	if err := run(destination, "git", "commit", "-m", "chore: initialize repository"); err != nil {
		_ = os.RemoveAll(destination)
		return Result{}, err
	}
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	_ = s.audit.Record(audit.Event{Operation: "repo:scaffold", Target: destination, Outcome: "success", Metadata: map[string]any{"files": len(paths)}})
	return Result{Path: destination, Files: paths}, nil
}

func (s *Service) Publish(ctx context.Context, path, repositoryName, visibility, confirmation string) (map[string]any, error) {
	if confirmation != security.Confirmation {
		return nil, fmt.Errorf("publishing requires confirmation=%q", security.Confirmation)
	}
	if !safeName.MatchString(repositoryName) {
		return nil, errors.New("invalid repository name")
	}
	if visibility != "private" && visibility != "public" && visibility != "internal" {
		return nil, errors.New("visibility must be private, public, or internal")
	}
	parent, err := s.containingBootstrapRoot(path)
	if err != nil {
		return nil, err
	}
	repoPath, err := security.ResolveExisting(parent, strings.TrimPrefix(path, parent+string(filepath.Separator)))
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(filepath.Join(repoPath, ".git")); err != nil {
		return nil, errors.New("path is not an initialized git repository")
	}
	ctx, cancel := context.WithTimeout(ctx, s.cfg.Timeout())
	defer cancel()
	slug := s.cfg.Organization + "/" + repositoryName
	if err := s.audit.Record(audit.Event{Operation: "repo:publish", Repository: slug, Target: repoPath, Outcome: "requested"}); err != nil {
		return nil, fmt.Errorf("write blocked because audit log is unavailable: %w", err)
	}
	cmd := exec.CommandContext(ctx, "gh", "repo", "create", slug, "--"+visibility, "--source", repoPath, "--remote", "origin", "--push")
	out, err := cmd.CombinedOutput()
	outcome := "success"
	if err != nil {
		outcome = "failure"
	}
	_ = s.audit.Record(audit.Event{Operation: "repo:publish", Repository: slug, Target: repoPath, Outcome: outcome})
	if err != nil {
		return nil, fmt.Errorf("gh repo create failed: %s", strings.TrimSpace(string(out)))
	}
	return map[string]any{"repository": slug, "path": repoPath, "output": strings.TrimSpace(string(out))}, nil
}

func (s *Service) allowedParent(requested string) (string, error) {
	requested, err := filepath.Abs(requested)
	if err != nil {
		return "", err
	}
	requested, err = filepath.EvalSymlinks(requested)
	if err != nil {
		return "", err
	}
	for _, root := range s.cfg.BootstrapRoots {
		rootReal, err := filepath.EvalSymlinks(root)
		if err == nil && requested == rootReal {
			return rootReal, nil
		}
	}
	return "", errors.New("parent is not an allowed bootstrap root")
}

func (s *Service) containingBootstrapRoot(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	for _, root := range s.cfg.BootstrapRoots {
		rel, err := filepath.Rel(root, abs)
		if err == nil && rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return root, nil
		}
	}
	return "", errors.New("path is outside allowed bootstrap roots")
}

func writeAtomically(destination string, files map[string]string) error {
	parent := filepath.Dir(destination)
	temp, err := os.MkdirTemp(parent, ".bifunctional-scaffold-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)
	for relative, content := range files {
		path := filepath.Join(temp, filepath.Clean(relative))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}
	return os.Rename(temp, destination)
}

func run(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s failed: %s", name, strings.TrimSpace(string(out)))
	}
	return nil
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = regexp.MustCompile(`[^a-z0-9._-]+`).ReplaceAllString(value, "-")
	return strings.Trim(value, "-.")
}

func validateOptions(req Request) error {
	switch strings.ToLower(req.License) {
	case "", "mit", "apache", "apache-2.0", "proprietary", "closed":
	default:
		return errors.New("license must be mit, apache-2.0, or proprietary")
	}
	switch strings.ToLower(req.Gitignore) {
	case "", "node", "go", "python":
	default:
		return errors.New("gitignore must be node, go, or python")
	}
	return nil
}

func projectFiles(req Request, brief, brand, design string) map[string]string {
	license := licenseText(req.License)
	return map[string]string{
		"README.md":                          readme(req.Project),
		"AGENTS.md":                          agents(req.Project),
		"CONTRIBUTING.md":                    contributing(),
		"SECURITY.md":                        securityPolicy(),
		"LICENSE":                            license,
		".gitignore":                         gitignore(req.Gitignore),
		".env.example":                       "# Document required environment variables without secret values.\n",
		"docs/index.md":                      docsIndex(req.Project),
		"docs/architecture.md":               architecture(req.Project),
		"docs/roadmap.md":                    roadmap(),
		"docs/design/brand-direction.md":     brand,
		"docs/design/design-system.md":       design,
		"docs/research/project-brief.md":     brief,
		"docs/decisions/README.md":           "# Architecture decisions\n\nRecord durable technical and product decisions here.\n",
		"docs/issues/README.md":              "# Issue knowledge\n\nUse this folder for durable issue summaries, not live status duplication.\n",
		"docs/notes/README.md":               "# Working notes\n\nPromote durable knowledge into canonical docs.\n",
		".github/pull_request_template.md":   pullRequestTemplate(),
		".github/ISSUE_TEMPLATE/feature.yml": featureIssueTemplate(),
	}
}

func readme(p ideation.Project) string {
	return fmt.Sprintf("# %s\n\n%s\n\n## Audience\n\n%s\n\n## Development\n\nDocument setup, build, test, and release commands as the stack is selected.\n\n## Knowledge\n\nStart with `docs/index.md`.\n", p.Name, p.Purpose, p.Audience)
}
func agents(p ideation.Project) string {
	return fmt.Sprintf("# Agent instructions\n\n## Intent\n\n%s\n\n## Operating rules\n\n- Read `docs/index.md` and relevant design documents before editing.\n- Keep changes scoped and preserve user work.\n- Never expose secrets or private knowledge in logs, prompts, issues, or generated files.\n- Add dependencies only for a demonstrated requirement.\n- Run the documented formatter, tests, and build before finishing.\n- Use conventional commits and include rationale in pull requests.\n- Update canonical docs when behavior or architecture changes.\n", p.Purpose)
}
func contributing() string {
	return "# Contributing\n\n1. Read `AGENTS.md` and `docs/index.md`.\n2. Create a focused branch and issue-linked change.\n3. Add or update tests with behavioral changes.\n4. Run formatting, tests, linting, and the production build.\n5. Open a pull request that explains intent, risk, verification, and rollback.\n"
}
func securityPolicy() string {
	return "# Security policy\n\nDo not commit credentials, personal data, private client material, or production exports. Report vulnerabilities privately to the repository maintainers. Rotate any credential that enters git history or logs. Treat generated content as untrusted input and review it before publication.\n"
}
func docsIndex(p ideation.Project) string {
	return fmt.Sprintf("---\ntitle: %s knowledge index\nowner: bifunctional\nstatus: active\nupdated: %s\ntags: [index, okf]\n---\n\n# Knowledge index\n\n- [Architecture](architecture.md)\n- [Roadmap](roadmap.md)\n- [Project brief](research/project-brief.md)\n- [Brand direction](design/brand-direction.md)\n- [Design system](design/design-system.md)\n- [Decisions](decisions/README.md)\n", p.Name, timeNow())
}
func architecture(p ideation.Project) string {
	return fmt.Sprintf("---\ntitle: %s architecture\nowner: bifunctional\nstatus: draft\nupdated: %s\ntags: [architecture, okf]\n---\n\n# Architecture\n\n## Context\n\n%s\n\n## Quality attributes\n\n- Security and privacy by default.\n- Observable failure behavior.\n- Reversible changes and documented recovery.\n- Accessibility and performance budgets appropriate to the product.\n\n## System design\n\nDefine components, data boundaries, dependencies, deployment, and failure modes before implementation hardens.\n", p.Name, timeNow(), p.Purpose)
}
func roadmap() string {
	return fmt.Sprintf("---\ntitle: Project roadmap\nowner: bifunctional\nstatus: draft\nupdated: %s\ntags: [roadmap, okf]\n---\n\n# Roadmap\n\n## Now\n\n- Validate the problem and first release boundary.\n- Establish architecture, design foundations, and delivery checks.\n\n## Next\n\n- Deliver and measure the first complete user workflow.\n\n## Later\n\n- Expand only from observed user and operational evidence.\n", timeNow())
}
func pullRequestTemplate() string {
	return "## Intent\n\n## Changes\n\n## Risk and rollback\n\n## Verification\n\n- [ ] Tests\n- [ ] Build\n- [ ] Docs\n- [ ] Security/privacy review\n"
}
func featureIssueTemplate() string {
	return "name: Feature\ndescription: Propose a scoped product or platform change\ntitle: \"feat: \"\nbody:\n  - type: textarea\n    id: outcome\n    attributes:\n      label: Desired outcome\n    validations:\n      required: true\n  - type: textarea\n    id: acceptance\n    attributes:\n      label: Acceptance criteria\n    validations:\n      required: true\n  - type: textarea\n    id: risks\n    attributes:\n      label: Risks and constraints\n"
}

func licenseText(kind string) string {
	switch strings.ToLower(kind) {
	case "apache-2.0", "apache":
		return apacheLicense
	case "proprietary", "closed":
		return "Copyright Bifunctional Labs Co. All rights reserved.\n\nThis source code and associated materials are proprietary and confidential. No license is granted except by written agreement.\n"
	default:
		return "MIT License\n\nCopyright (c) Bifunctional Labs Co\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all\ncopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\nSOFTWARE.\n"
	}
}

func gitignore(kind string) string {
	base := ".DS_Store\n.env\n.env.*\n!.env.example\n*.log\ncoverage/\ndist/\nbuild/\n"
	switch strings.ToLower(kind) {
	case "go":
		return base + "bin/\n*.test\n"
	case "python":
		return base + ".venv/\n__pycache__/\n*.py[cod]\n.pytest_cache/\n"
	default:
		return base + "node_modules/\n.astro/\n.next/\n"
	}
}

func timeNow() string { return now().Format("2006-01-02") }

var now = func() time.Time { return time.Now().UTC() }
