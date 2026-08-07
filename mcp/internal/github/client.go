package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/audit"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/config"
	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/security"
)

type Client struct {
	cfg   *config.Config
	audit *audit.Logger
}

type IssueCreate struct {
	Repository   string
	Title        string
	Body         string
	Labels       []string
	Assignees    []string
	Milestone    string
	Confirmation string
}

type IssueUpdate struct {
	Repository   string
	Number       int
	Title        string
	Body         string
	State        string
	AddLabels    []string
	Confirmation string
}

type PullRequestCreate struct {
	Repository   string
	Title        string
	Body         string
	Head         string
	Base         string
	Draft        bool
	Confirmation string
}

var branchPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/-]{0,127}$`)

func New(cfg *config.Config, logger *audit.Logger) *Client {
	return &Client{cfg: cfg, audit: logger}
}

func (c *Client) SearchIssues(ctx context.Context, repository, query, state string, limit int) (any, error) {
	repo, err := c.authorize(repository, "issues:read", "")
	if err != nil {
		return nil, err
	}
	if state == "" {
		state = "open"
	}
	if state != "open" && state != "closed" && state != "all" {
		return nil, errors.New("state must be open, closed, or all")
	}
	if limit <= 0 || limit > 100 {
		limit = 30
	}
	args := []string{"issue", "list", "--repo", repo.Slug, "--state", state, "--limit", strconv.Itoa(limit), "--json", "number,title,state,url,labels,assignees,milestone,updatedAt"}
	if query != "" {
		args = append(args, "--search", query)
	}
	return c.runJSON(ctx, args...)
}

func (c *Client) Labels(ctx context.Context, repository string, limit int) (any, error) {
	repo, err := c.authorize(repository, "issues:read", "")
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	return c.runJSON(ctx, "label", "list", "--repo", repo.Slug, "--limit", strconv.Itoa(limit), "--json", "name,color,description")
}

func (c *Client) Milestones(ctx context.Context, repository, state string) (any, error) {
	repo, err := c.authorize(repository, "issues:read", "")
	if err != nil {
		return nil, err
	}
	if state == "" {
		state = "open"
	}
	if state != "open" && state != "closed" && state != "all" {
		return nil, errors.New("state must be open, closed, or all")
	}
	return c.runJSON(ctx, "api", "--method", "GET", "repos/"+repo.Slug+"/milestones", "-f", "state="+state, "-f", "per_page=100")
}

func (c *Client) Repository(ctx context.Context, repository string) (any, error) {
	repo, err := c.authorize(repository, "repo:read", "")
	if err != nil {
		return nil, err
	}
	return c.runJSON(ctx, "repo", "view", repo.Slug, "--json", "nameWithOwner,description,visibility,defaultBranchRef,url,isArchived,isFork,licenseInfo,repositoryTopics")
}

func (c *Client) CreateIssue(ctx context.Context, input IssueCreate) (any, error) {
	repo, err := c.authorize(input.Repository, "issues:write", input.Confirmation)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.Title) == "" {
		return nil, errors.New("title is required")
	}
	if err := c.begin("issues:create", repo.Slug, input.Title); err != nil {
		return nil, err
	}
	args := []string{"issue", "create", "--repo", repo.Slug, "--title", input.Title, "--body", input.Body}
	for _, label := range input.Labels {
		args = append(args, "--label", label)
	}
	for _, assignee := range input.Assignees {
		args = append(args, "--assignee", assignee)
	}
	if input.Milestone != "" {
		args = append(args, "--milestone", input.Milestone)
	}
	output, err := c.runText(ctx, args...)
	c.record("issues:create", repo.Slug, input.Title, err)
	if err != nil {
		return nil, err
	}
	return map[string]any{"url": strings.TrimSpace(output)}, nil
}

func (c *Client) UpdateIssue(ctx context.Context, input IssueUpdate) (any, error) {
	repo, err := c.authorize(input.Repository, "issues:write", input.Confirmation)
	if err != nil {
		return nil, err
	}
	if input.Number <= 0 {
		return nil, errors.New("number must be positive")
	}
	if err := c.begin("issues:update", repo.Slug, strconv.Itoa(input.Number)); err != nil {
		return nil, err
	}
	args := []string{"issue", "edit", strconv.Itoa(input.Number), "--repo", repo.Slug}
	if input.Title != "" {
		args = append(args, "--title", input.Title)
	}
	if input.Body != "" {
		args = append(args, "--body", input.Body)
	}
	for _, label := range input.AddLabels {
		args = append(args, "--add-label", label)
	}
	if len(args) == 5 && input.State == "" {
		return nil, errors.New("no issue fields supplied")
	}
	if len(args) > 5 {
		_, err = c.runText(ctx, args...)
	}
	if err == nil && input.State != "" {
		if input.State != "open" && input.State != "closed" {
			return nil, errors.New("state must be open or closed")
		}
		action := "reopen"
		if input.State == "closed" {
			action = "close"
		}
		_, err = c.runText(ctx, "issue", action, strconv.Itoa(input.Number), "--repo", repo.Slug)
	}
	c.record("issues:update", repo.Slug, strconv.Itoa(input.Number), err)
	if err != nil {
		return nil, err
	}
	return c.runJSON(ctx, "issue", "view", strconv.Itoa(input.Number), "--repo", repo.Slug, "--json", "number,title,state,url,labels,assignees,milestone,updatedAt")
}

func (c *Client) CreatePullRequest(ctx context.Context, input PullRequestCreate) (any, error) {
	repo, err := c.authorize(input.Repository, "pulls:write", input.Confirmation)
	if err != nil {
		return nil, err
	}
	if input.Title == "" || input.Head == "" {
		return nil, errors.New("title and head are required")
	}
	if err := c.begin("pulls:create", repo.Slug, input.Head); err != nil {
		return nil, err
	}
	if input.Base == "" {
		input.Base = repo.DefaultBranch
	}
	args := []string{"pr", "create", "--repo", repo.Slug, "--title", input.Title, "--body", input.Body, "--head", input.Head, "--base", input.Base}
	if input.Draft {
		args = append(args, "--draft")
	}
	out, err := c.runText(ctx, args...)
	c.record("pulls:create", repo.Slug, input.Head, err)
	if err != nil {
		return nil, err
	}
	return map[string]any{"url": strings.TrimSpace(out)}, nil
}

func (c *Client) CreateBranch(ctx context.Context, repository, branch, base, confirmation string) (any, error) {
	repo, err := c.authorize(repository, "branches:write", confirmation)
	if err != nil {
		return nil, err
	}
	if !branchPattern.MatchString(branch) || strings.Contains(branch, "..") || strings.HasSuffix(branch, "/") {
		return nil, errors.New("invalid branch name")
	}
	if base == "" {
		base = repo.DefaultBranch
	}
	if !branchPattern.MatchString(base) || strings.Contains(base, "..") {
		return nil, errors.New("invalid base branch")
	}
	if err := c.begin("branches:create", repo.Slug, branch); err != nil {
		return nil, err
	}
	value, err := c.runJSON(ctx, "api", "repos/"+repo.Slug+"/git/ref/heads/"+base)
	if err != nil {
		c.record("branches:create", repo.Slug, branch, err)
		return nil, err
	}
	root, ok := value.(map[string]any)
	if !ok {
		return nil, errors.New("unexpected branch response")
	}
	object, ok := root["object"].(map[string]any)
	if !ok {
		return nil, errors.New("branch response has no object")
	}
	sha, ok := object["sha"].(string)
	if !ok || sha == "" {
		return nil, errors.New("branch response has no commit SHA")
	}
	out, err := c.runJSON(ctx, "api", "--method", "POST", "repos/"+repo.Slug+"/git/refs", "-f", "ref=refs/heads/"+branch, "-f", "sha="+sha)
	c.record("branches:create", repo.Slug, branch, err)
	return out, err
}

func (c *Client) ReleaseNotes(ctx context.Context, repository, query string, limit int) (any, error) {
	repo, err := c.authorize(repository, "pulls:read", "")
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	args := []string{"pr", "list", "--repo", repo.Slug, "--state", "merged", "--limit", strconv.Itoa(limit), "--json", "number,title,url,author,labels,mergedAt,baseRefName,headRefName"}
	if query != "" {
		args = append(args, "--search", query)
	}
	value, err := c.runJSON(ctx, args...)
	if err != nil {
		return nil, err
	}
	items, _ := value.([]any)
	var notes strings.Builder
	notes.WriteString("# Release notes\n\n")
	if len(items) == 0 {
		notes.WriteString("No merged pull requests matched the requested range.\n")
	}
	for _, item := range items {
		pr, _ := item.(map[string]any)
		number, _ := pr["number"].(float64)
		title, _ := pr["title"].(string)
		url, _ := pr["url"].(string)
		fmt.Fprintf(&notes, "- [#%d %s](%s)\n", int(number), title, url)
	}
	return map[string]any{"markdown": notes.String(), "pullRequests": value}, nil
}

func (c *Client) authorize(slug, operation, confirmation string) (config.Repository, error) {
	repo, ok := c.cfg.Repository(slug)
	if !ok {
		return config.Repository{}, fmt.Errorf("repository %q is not registered", slug)
	}
	if err := security.Require(repo, operation, confirmation); err != nil {
		return config.Repository{}, err
	}
	return repo, nil
}

func (c *Client) runJSON(ctx context.Context, args ...string) (any, error) {
	out, err := c.runText(ctx, args...)
	if err != nil {
		return nil, err
	}
	var value any
	if err := json.Unmarshal([]byte(out), &value); err != nil {
		return nil, fmt.Errorf("decode gh response: %w", err)
	}
	return value, nil
}

func (c *Client) runText(ctx context.Context, args ...string) (string, error) {
	commandCtx, cancel := context.WithTimeout(ctx, c.cfg.Timeout())
	defer cancel()
	cmd := exec.CommandContext(commandCtx, "gh", args...)
	out, err := cmd.CombinedOutput()
	if commandCtx.Err() != nil {
		return "", fmt.Errorf("gh command timed out after %s", c.cfg.Timeout())
	}
	if err != nil {
		return "", fmt.Errorf("gh command failed: %s", strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func (c *Client) record(operation, repository, target string, err error) {
	outcome := "success"
	if err != nil {
		outcome = "failure"
	}
	_ = c.audit.Record(audit.Event{Operation: operation, Repository: repository, Target: target, Outcome: outcome, Metadata: map[string]any{"at": time.Now().UTC()}})
}

func (c *Client) begin(operation, repository, target string) error {
	if err := c.audit.Record(audit.Event{Operation: operation, Repository: repository, Target: target, Outcome: "requested"}); err != nil {
		return fmt.Errorf("write blocked because audit log is unavailable: %w", err)
	}
	return nil
}
