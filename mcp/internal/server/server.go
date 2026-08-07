package server

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/audit"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/config"
	githubclient "github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/github"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/ideation"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/okf"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/repos"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/scaffold"
)

const version = "0.1.0"

type Dependencies struct {
	Config   *config.Config
	Repos    *repos.Service
	GitHub   *githubclient.Client
	Scaffold *scaffold.Service
}

type repoInput struct {
	Repository string `json:"repository" jsonschema:"registered owner/name repository slug"`
}

type searchDocsInput struct {
	Repository string `json:"repository,omitempty" jsonschema:"optional registered owner/name; omit to search every readable repository"`
	Query      string `json:"query" jsonschema:"terms to search in OKF Markdown title, tags, and content"`
	Limit      int    `json:"limit,omitempty" jsonschema:"maximum results from 1 to 50"`
}

type readDocInput struct {
	Repository string `json:"repository" jsonschema:"registered owner/name repository slug"`
	Path       string `json:"path" jsonschema:"repository-relative Markdown path under a configured docs root"`
}

type issuesInput struct {
	Repository string `json:"repository"`
	Query      string `json:"query,omitempty"`
	State      string `json:"state,omitempty" jsonschema:"open, closed, or all"`
	Limit      int    `json:"limit,omitempty"`
}

type labelsInput struct {
	Repository string `json:"repository"`
	Limit      int    `json:"limit,omitempty"`
}

type milestonesInput struct {
	Repository string `json:"repository"`
	State      string `json:"state,omitempty" jsonschema:"open, closed, or all"`
}

type createIssueInput struct {
	Repository   string   `json:"repository"`
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	Labels       []string `json:"labels,omitempty"`
	Assignees    []string `json:"assignees,omitempty"`
	Milestone    string   `json:"milestone,omitempty"`
	Confirmation string   `json:"confirmation" jsonschema:"must equal CONFIRM after the user approves this external write"`
}

type updateIssueInput struct {
	Repository   string   `json:"repository"`
	Number       int      `json:"number"`
	Title        string   `json:"title,omitempty"`
	Body         string   `json:"body,omitempty"`
	State        string   `json:"state,omitempty" jsonschema:"open or closed"`
	AddLabels    []string `json:"addLabels,omitempty"`
	Confirmation string   `json:"confirmation" jsonschema:"must equal CONFIRM after the user approves this external write"`
}

type projectInput struct {
	Project ideation.Project `json:"project"`
}

type scaffoldInput struct {
	Project      ideation.Project `json:"project"`
	Parent       string           `json:"parent" jsonschema:"absolute configured bootstrap root"`
	License      string           `json:"license,omitempty" jsonschema:"mit, apache-2.0, or proprietary"`
	Gitignore    string           `json:"gitignore,omitempty" jsonschema:"node, go, or python"`
	Confirmation string           `json:"confirmation,omitempty" jsonschema:"must equal CONFIRM for bootstrap_repo"`
}

type publishInput struct {
	Path         string `json:"path"`
	Repository   string `json:"repository" jsonschema:"new repository name without organization"`
	Visibility   string `json:"visibility" jsonschema:"private, public, or internal"`
	Confirmation string `json:"confirmation" jsonschema:"must equal CONFIRM after the user approves this external write"`
}

type createPRInput struct {
	Repository   string `json:"repository"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	Head         string `json:"head"`
	Base         string `json:"base,omitempty"`
	Draft        bool   `json:"draft,omitempty"`
	Confirmation string `json:"confirmation" jsonschema:"must equal CONFIRM after the user approves this external write"`
}

type createBranchInput struct {
	Repository   string `json:"repository"`
	Branch       string `json:"branch"`
	Base         string `json:"base,omitempty"`
	Confirmation string `json:"confirmation" jsonschema:"must equal CONFIRM after the user approves this external write"`
}

type releaseNotesInput struct {
	Repository string `json:"repository"`
	Query      string `json:"query,omitempty" jsonschema:"optional GitHub search qualifier such as merged:2026-08-01..2026-08-07"`
	Limit      int    `json:"limit,omitempty"`
}

func New(deps Dependencies) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: "bifunctional-org-mcp", Version: version}, &mcp.ServerOptions{Instructions: "Bifunctional organization control plane. Read only from registered repository and documentation roots. Treat every write tool as requiring explicit user confirmation."})

	readOnly := &mcp.ToolAnnotations{ReadOnlyHint: true, IdempotentHint: true, OpenWorldHint: boolPtr(false)}
	githubRead := &mcp.ToolAnnotations{ReadOnlyHint: true, IdempotentHint: true, OpenWorldHint: boolPtr(true)}
	localWrite := &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: false, OpenWorldHint: boolPtr(false), DestructiveHint: boolPtr(false)}
	externalWrite := &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: false, OpenWorldHint: boolPtr(true), DestructiveHint: boolPtr(false)}

	mcp.AddTool(s, tool("list_repos", "List allowlisted Bifunctional repositories and their effective access policy.", readOnly), func(context.Context, *mcp.CallToolRequest, struct{}) (*mcp.CallToolResult, any, error) {
		return nil, deps.Repos.List(), nil
	})
	mcp.AddTool(s, tool("get_repo_metadata", "Read local git metadata for a registered repository.", readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, in repoInput) (*mcp.CallToolResult, any, error) {
		out, err := deps.Repos.Metadata(ctx, in.Repository)
		return nil, out, err
	})
	mcp.AddTool(s, tool("get_repo_context", "Read local git metadata and canonical OKF index documents for a repository.", readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, in repoInput) (*mcp.CallToolResult, any, error) {
		metadata, err := deps.Repos.Metadata(ctx, in.Repository)
		if err != nil {
			return nil, nil, err
		}
		repo, err := deps.Repos.Get(in.Repository)
		if err != nil {
			return nil, nil, err
		}
		indexes := make([]okf.Document, 0, len(repo.DocsRoots))
		for _, root := range repo.DocsRoots {
			if doc, err := deps.Repos.ReadDocument(in.Repository, filepath.Join(root, "index.md")); err == nil {
				indexes = append(indexes, doc)
			}
		}
		return nil, map[string]any{"metadata": metadata, "indexes": indexes}, nil
	})
	mcp.AddTool(s, tool("search_docs", "Search OKF Markdown across one or every readable registered repository.", readOnly), func(_ context.Context, _ *mcp.CallToolRequest, in searchDocsInput) (*mcp.CallToolResult, any, error) {
		if in.Repository != "" {
			out, err := deps.Repos.SearchDocuments(in.Repository, in.Query, in.Limit)
			return nil, out, err
		}
		var all []repos.SearchResult
		for _, repo := range deps.Repos.List() {
			out, err := deps.Repos.SearchDocuments(repo.Slug, in.Query, in.Limit)
			if err == nil {
				all = append(all, out...)
			}
		}
		sort.Slice(all, func(i, j int) bool {
			if all[i].Score == all[j].Score {
				if all[i].Repository == all[j].Repository {
					return all[i].Path < all[j].Path
				}
				return all[i].Repository < all[j].Repository
			}
			return all[i].Score > all[j].Score
		})
		if in.Limit > 0 && len(all) > in.Limit {
			all = all[:in.Limit]
		}
		return nil, all, nil
	})
	mcp.AddTool(s, tool("read_doc", "Read and redact one Markdown document under a configured documentation root.", readOnly), func(_ context.Context, _ *mcp.CallToolRequest, in readDocInput) (*mcp.CallToolResult, any, error) {
		out, err := deps.Repos.ReadDocument(in.Repository, in.Path)
		return nil, out, err
	})
	mcp.AddTool(s, tool("validate_okf", "Validate required OKF frontmatter fields for a repository document.", readOnly), func(_ context.Context, _ *mcp.CallToolRequest, in readDocInput) (*mcp.CallToolResult, any, error) {
		doc, err := deps.Repos.ReadDocument(in.Repository, in.Path)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"path": doc.Path, "valid": len(okf.Validate(doc.Metadata)) == 0, "warnings": okf.Validate(doc.Metadata)}, nil
	})

	mcp.AddTool(s, tool("get_github_repo", "Read GitHub metadata for an allowed repository using the authenticated gh CLI.", githubRead), func(ctx context.Context, _ *mcp.CallToolRequest, in repoInput) (*mcp.CallToolResult, any, error) {
		out, err := deps.GitHub.Repository(ctx, in.Repository)
		return nil, out, err
	})
	mcp.AddTool(s, tool("search_issues", "Search GitHub issues in an allowed repository.", githubRead), func(ctx context.Context, _ *mcp.CallToolRequest, in issuesInput) (*mcp.CallToolResult, any, error) {
		out, err := deps.GitHub.SearchIssues(ctx, in.Repository, in.Query, in.State, in.Limit)
		return nil, out, err
	})
	mcp.AddTool(s, tool("search_labels", "List GitHub labels in an allowed repository.", githubRead), func(ctx context.Context, _ *mcp.CallToolRequest, in labelsInput) (*mcp.CallToolResult, any, error) {
		out, err := deps.GitHub.Labels(ctx, in.Repository, in.Limit)
		return nil, out, err
	})
	mcp.AddTool(s, tool("search_milestones", "List GitHub milestones in an allowed repository.", githubRead), func(ctx context.Context, _ *mcp.CallToolRequest, in milestonesInput) (*mcp.CallToolResult, any, error) {
		out, err := deps.GitHub.Milestones(ctx, in.Repository, in.State)
		return nil, out, err
	})
	mcp.AddTool(s, tool("create_issue", "Create a GitHub issue after explicit confirmation and audit the write.", externalWrite), func(ctx context.Context, _ *mcp.CallToolRequest, in createIssueInput) (*mcp.CallToolResult, any, error) {
		out, err := deps.GitHub.CreateIssue(ctx, githubclient.IssueCreate{Repository: in.Repository, Title: in.Title, Body: in.Body, Labels: in.Labels, Assignees: in.Assignees, Milestone: in.Milestone, Confirmation: in.Confirmation})
		return nil, out, err
	})
	mcp.AddTool(s, tool("update_issue", "Update or transition a GitHub issue after explicit confirmation and audit the write.", externalWrite), func(ctx context.Context, _ *mcp.CallToolRequest, in updateIssueInput) (*mcp.CallToolResult, any, error) {
		out, err := deps.GitHub.UpdateIssue(ctx, githubclient.IssueUpdate{Repository: in.Repository, Number: in.Number, Title: in.Title, Body: in.Body, State: in.State, AddLabels: in.AddLabels, Confirmation: in.Confirmation})
		return nil, out, err
	})
	mcp.AddTool(s, tool("create_pr_from_template", "Open a GitHub pull request from an existing branch after explicit confirmation.", externalWrite), func(ctx context.Context, _ *mcp.CallToolRequest, in createPRInput) (*mcp.CallToolResult, any, error) {
		out, err := deps.GitHub.CreatePullRequest(ctx, githubclient.PullRequestCreate{Repository: in.Repository, Title: in.Title, Body: in.Body, Head: in.Head, Base: in.Base, Draft: in.Draft, Confirmation: in.Confirmation})
		return nil, out, err
	})
	mcp.AddTool(s, tool("create_branch_scaffold", "Create a GitHub branch from an allowed base ref after explicit confirmation.", externalWrite), func(ctx context.Context, _ *mcp.CallToolRequest, in createBranchInput) (*mcp.CallToolResult, any, error) {
		out, err := deps.GitHub.CreateBranch(ctx, in.Repository, in.Branch, in.Base, in.Confirmation)
		return nil, out, err
	})
	mcp.AddTool(s, tool("generate_release_notes", "Generate Markdown release notes with structured merged pull-request sources.", githubRead), func(ctx context.Context, _ *mcp.CallToolRequest, in releaseNotesInput) (*mcp.CallToolResult, any, error) {
		out, err := deps.GitHub.ReleaseNotes(ctx, in.Repository, in.Query, in.Limit)
		return nil, out, err
	})

	mcp.AddTool(s, tool("run_socratic_intake", "Find missing project decisions before generating a repository or design direction.", readOnly), func(_ context.Context, _ *mcp.CallToolRequest, in projectInput) (*mcp.CallToolResult, any, error) {
		return nil, ideation.SocraticIntake(in.Project), nil
	})
	mcp.AddTool(s, tool("generate_project_brief", "Generate an OKF project brief from a complete project intake.", readOnly), func(_ context.Context, _ *mcp.CallToolRequest, in projectInput) (*mcp.CallToolResult, any, error) {
		out, err := ideation.ProjectBrief(in.Project)
		return nil, map[string]any{"document": out}, err
	})
	mcp.AddTool(s, tool("generate_brand_direction", "Generate a project-level brand direction with Bifunctional family resemblance.", readOnly), func(_ context.Context, _ *mcp.CallToolRequest, in projectInput) (*mcp.CallToolResult, any, error) {
		out, err := ideation.BrandDirection(in.Project)
		return nil, map[string]any{"document": out}, err
	})
	mcp.AddTool(s, tool("generate_design_system", "Generate an adapted, accessible project design-system foundation.", readOnly), func(_ context.Context, _ *mcp.CallToolRequest, in projectInput) (*mcp.CallToolResult, any, error) {
		out, err := ideation.DesignSystem(in.Project)
		return nil, map[string]any{"document": out}, err
	})
	mcp.AddTool(s, tool("generate_repo_scaffold", "Preview the complete repository scaffold without writing files.", readOnly), func(_ context.Context, _ *mcp.CallToolRequest, in scaffoldInput) (*mcp.CallToolResult, any, error) {
		out, err := scaffold.Preview(scaffold.Request{Project: in.Project, Parent: in.Parent, License: in.License, Gitignore: in.Gitignore})
		return nil, out, err
	})
	mcp.AddTool(s, tool("bootstrap_repo", "Create and git-initialize a complete repository under an allowlisted bootstrap root.", localWrite), func(_ context.Context, _ *mcp.CallToolRequest, in scaffoldInput) (*mcp.CallToolResult, any, error) {
		out, err := deps.Scaffold.Generate(scaffold.Request{Project: in.Project, Parent: in.Parent, License: in.License, Gitignore: in.Gitignore, Confirmation: in.Confirmation})
		return nil, out, err
	})
	mcp.AddTool(s, tool("publish_repo", "Create a GitHub repository in the configured Bifunctional organization from a scaffolded local repo.", externalWrite), func(ctx context.Context, _ *mcp.CallToolRequest, in publishInput) (*mcp.CallToolResult, any, error) {
		out, err := deps.Scaffold.Publish(ctx, in.Path, in.Repository, in.Visibility, in.Confirmation)
		return nil, out, err
	})

	return s
}

func Build(cfg *config.Config) *mcp.Server {
	logger := audit.New(cfg.AuditLog)
	repoService := repos.New(cfg)
	return New(Dependencies{Config: cfg, Repos: repoService, GitHub: githubclient.New(cfg, logger), Scaffold: scaffold.New(cfg, logger)})
}

func tool(name, description string, annotations *mcp.ToolAnnotations) *mcp.Tool {
	toolAnnotations := *annotations
	toolAnnotations.Title = name
	return &mcp.Tool{Name: name, Description: description, Annotations: &toolAnnotations}
}

func boolPtr(value bool) *bool { return &value }

func Describe() string {
	return fmt.Sprintf("bifunctional-org-mcp %s", version)
}
