# Central Org MCP Architecture Spike

Status: draft
Scope: Bifunctional org repo space and local MCP control plane

## Executive Summary

The right shape for this system is a small central control plane plus many repo-local knowledge folders.

The control plane should:

- discover local clones and known GitHub repos
- read and search repo docs safely
- surface GitHub issues, labels, milestones, and repo metadata
- guide new project ideation through a Socratic intake
- bootstrap new repos with a consistent Bifunctional package
- keep public brand material separate from sensitive internal knowledge

The repo-local knowledge layer should:

- live in each repo as a `docs/` folder or equivalent
- use a simple OKF-style markdown convention
- be readable by both humans and the MCP server
- remain independent from the public website build

The safest version is read-first, write-second.
Start with indexing and query, then add repo creation and cross-repo mutations once the access model is proven.

## What We Are Building

The user goal is not just "an MCP server".
It is a central org brain for Bifunctional that can:

1. understand every repo in the org space
2. answer questions from repo knowledge
3. manage work at the GitHub issue level
4. help create new projects with a brand-consistent intake
5. generate a complete repo starter kit

This is effectively three layers:

- `control plane`: org registry, policy, routing, and permissions
- `knowledge plane`: repo docs, design docs, decisions, issue knowledge
- `execution plane`: repo creation, issue updates, branch scaffolding, and optional write actions

## Design Principles

### 1. Local first

Prefer local git clones already on the machine.
That keeps the system fast, private, and usable even when GitHub API access is limited.

### 2. Read by default

The first release should be able to read and summarize.
Write actions should be opt-in and scoped.

### 3. Knowledge stays close to code

Every repo should carry its own docs folder.
The MCP server should not become the only place where project context lives.

### 4. Public and sensitive content must stay separated

Public brand pages can be indexed for context, but sensitive internal material should never be exported by default.

### 5. Small, explicit tools

Prefer narrow tools with clear names over one huge "do everything" endpoint.

## Recommended Repo Model

### Central org repo

Use one central repo as the control plane and reference hub.
It can hold:

- the repo registry
- the docs standard
- MCP server config
- shared agent instructions
- brand system references
- project templates

If this website repo is not meant to be that control plane, create a separate org admin repo.

### Repo-local docs

Each project repo should contain a knowledge wiki, ideally:

```text
docs/
  index.md
  architecture.md
  decisions/
  design/
  research/
  issues/
  roadmap.md
  notes/
```

This is the minimum useful structure.
It can be expanded later, but the first pass should stay lightweight.

### OKF-style convention

If OKF is the intended internal format, make it simple and consistent:

- markdown first
- stable headings
- frontmatter for metadata
- one canonical index per folder
- cross-links instead of duplication

Example frontmatter:

```md
---
title: Project Overview
owner: bifunctional
status: active
updated: 2026-08-07
tags: [mcp, docs, org]
---
```

## MCP Tool Set

### V1 read tools

- `list_repos`
- `get_repo_context`
- `search_docs`
- `read_doc`
- `search_issues`
- `search_labels`
- `search_milestones`
- `get_repo_metadata`

### V1 write tools

- `create_issue`
- `update_issue`
- `bootstrap_repo`
- `generate_repo_scaffold`

### V1 ideation tools

- `generate_project_brief`
- `generate_brand_direction`
- `generate_design_system`
- `run_socratic_intake`

### V2 tools

- `sync_docs_from_repo`
- `mirror_issue_state`
- `generate_release_notes`
- `create_branch_scaffold`
- `create_pr_from_template`

## Repo Registry

The MCP server needs a machine-readable registry of repos.

Suggested fields:

- repo slug
- local path
- visibility
- default branch
- docs root
- issue source
- sensitivity class
- allowed operations
- sync state

Example:

```json
{
  "repo": "BifunctionalLabsCo/bifunctionallabsco.github.io",
  "localPath": "/Users/zubinj/forge/bifunctionallabsco.github.io",
  "visibility": "public",
  "docsRoot": "colabspace",
  "issueSource": "github",
  "sensitivity": "mixed",
  "allowedOperations": ["read", "issue-read", "issue-write"]
}
```

## Access Model

The security model should be explicit.

### Read scopes

- repo metadata
- docs content
- issues, labels, milestones
- branch names

### Write scopes

- issue creation and updates
- repo bootstrap
- doc creation in allowed repos
- branch scaffolding

### Sensitive content policy

Do not index or export sensitive content broadly.

Recommended approach:

- classify repos by sensitivity
- redact secrets before indexing
- keep private repos queryable only when explicitly allowed
- never expose credentials, tokens, or personal data in model-facing summaries
- separate public website content from internal org knowledge

## Project Ideation Workflow

This is the Socratic part.

The server should help create a new project by asking a compact sequence of questions:

1. What is the project for?
2. Who is it for?
3. What kind of app or site is it?
4. What does success look like?
5. What should feel inherited from Bifunctional?
6. What should feel distinct to this project?
7. What are the first screens, pages, or workflows?
8. What assets, docs, or constraints already exist?

The output should include:

- project narrative
- visual direction
- adapted design system
- content tone
- repo scaffold plan
- initial docs plan

### Brand inheritance model

The project should inherit the core Bifunctional system:

- clarity over noise
- boutique over bloated
- structured but warm
- technical but human
- disciplined typography and spacing

Then adapt per project:

- palette
- imagery
- component emphasis
- tone of voice
- interaction style

This is not a clone of the website brand.
It is a family resemblance with room for the project to breathe.

## New Repo Bootstrap

When creating a new repo, the MCP server should scaffold a complete starter package.

Minimum files:

- `README.md`
- `CONTRIBUTING.md`
- `agents.md`
- `.gitignore`
- `LICENSE` when appropriate
- `docs/index.md`
- `docs/architecture.md`
- `docs/design/brand-direction.md`
- `docs/design/design-system.md`
- `docs/research/project-brief.md`
- `docs/issues/`

Recommended additions:

- `CHANGELOG.md`
- `docs/decisions/`
- `docs/roadmap.md`
- `docs/notes/`
- repo-specific `.env.example`

### `agents.md`

This should be the repo's operating manual for AI agents.
It should explain:

- project intent
- coding conventions
- docs conventions
- build/test commands
- safety rules
- branch and commit conventions

### `CONTRIBUTING.md`

This should be shared across the org where possible.
Keep it concise and enforce the same basics:

- read the brand and repo context first
- keep changes scoped
- do not add unnecessary dependencies
- run the build or checks before finishing

## GitHub Issue, Label, and Milestone Access

The central MCP server should be able to:

- search issues by repo
- filter by label or milestone
- create and update issues
- read cross-repo issue state
- optionally mirror issue summaries into docs

That gives the central org repo a single place to think without making GitHub the only knowledge store.

## Public Website Handling

The public Bifunctional website should be treated as public context, not a dumping ground for private knowledge.

Recommended rules:

- public website docs may be indexed
- internal strategy docs should stay private
- brand assets can be public only if intended
- no private client or internal repo content should be mirrored into the site repo

In practice:

- public brand material can inform the model
- sensitive internal material should be omitted or summarized with redaction

## What Not To Do

- do not make the MCP server a single giant opaque agent
- do not index every file blindly
- do not allow write access before the registry and permissions model exist
- do not duplicate the same long-form knowledge in several repos
- do not turn the public website repo into a secret management layer

## Recommended Implementation Order

### Phase 1

- define repo registry format
- define docs standard
- build read-only doc search
- build repo metadata lookup
- build GitHub issue and milestone read tools

### Phase 2

- add Socratic project intake
- generate project briefs
- generate brand direction
- generate design system docs
- add repo bootstrap scaffolding

### Phase 3

- add write actions for issues and docs
- add repo creation helpers
- add branch scaffolding
- add PR drafting flow

### Phase 4

- add multi-repo sync
- add release-note synthesis
- add issue-to-doc mirror workflows

## Proposed Central Repo Structure

If this repo becomes the control plane, a clean layout could be:

```text
colabspace/
  org/
    repo-registry.md
    docs-standard.md
    mcp-roadmap.md
    security-policy.md
  projects/
    active/
    archived/
  templates/
    repo-bootstrap/
    issue-templates/
    design/
```

This keeps the organization-level material separate from the public website content.

## Open Questions

These should be answered before implementation hardens:

- Is OKF already defined, or do we need to formalize it?
- Which repos are in scope on day one?
- Which repos are read-only versus writable?
- Do we want the central repo to be the website repo, or a separate org admin repo?
- What is the exact sensitivity classification model?
- Should issue state be mirrored into docs, or kept only in GitHub?
- Which write actions require manual confirmation?
- Do we want to support private repo indexing on the local machine only, or also remote indexing?

## Suggested First Deliverable

The first deliverable should be a docs-only spike package:

- the architecture spec
- the repo registry schema
- the OKF docs standard
- the bootstrap template
- the security policy

That is enough to start implementing the server without risking accidental exposure or a messy repo layout.
