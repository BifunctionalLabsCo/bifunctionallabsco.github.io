# Bifunctional organization MCP

A local-first MCP server for querying Bifunctional repositories, searching canonical OKF documentation, operating GitHub issues and pull requests, guiding project definition, and creating production-ready repository scaffolds.

The server is written in Go and uses the official Model Context Protocol SDK. It runs over stdio so credentials remain in the operator environment and no network listener is exposed by default.

The canonical organization hosting guidance, including the small-VPS posture and provider evaluation criteria, is in [`../colabspace/org/mcp-hosting.md`](../colabspace/org/mcp-hosting.md).

## Trust model

- A repository must exist in `config/repos.json` before it is visible.
- Every operation must be explicitly allowed per repository.
- Documentation reads are restricted to configured roots and Markdown files.
- Real paths are checked after symlink resolution.
- known secret files are denied and secret-like content is redacted.
- GitHub writes and local scaffolding require `confirmation: "CONFIRM"`.
- Writes append a JSON event to a mode `0600` audit log.
- The public Astro build does not ingest `mcp/` or `colabspace/`.

The registry is policy, not merely discovery configuration. Review changes to it as security changes.

## Tool surface

Repository and knowledge:

- `list_repos`
- `get_repo_metadata`
- `get_repo_context`
- `search_docs`
- `read_doc`
- `validate_okf`

GitHub:

- `get_github_repo`
- `search_issues`
- `search_labels`
- `search_milestones`
- `create_issue`
- `update_issue`
- `create_pr_from_template`
- `create_branch_scaffold`
- `generate_release_notes`

Project factory:

- `run_socratic_intake`
- `generate_project_brief`
- `generate_brand_direction`
- `generate_design_system`
- `generate_repo_scaffold`
- `bootstrap_repo`
- `publish_repo`

## Requirements

- Go 1.26 or newer for source builds.
- `git` for local repository metadata and bootstrapping.
- authenticated `gh` for GitHub tools.

Check GitHub access without exposing credentials:

```bash
gh auth status
```

## Build and verify

```bash
cd mcp
make check
```

## Run

```bash
cd mcp
make run
```

Use a private local registry for private repos:

```bash
cp config/repos.json config/repos.local.json
$EDITOR config/repos.local.json
CONFIG=config/repos.local.json make run
```

`config/repos.local.json` is ignored. Never add local private-repository paths or policies to the public registry unless their disclosure is intentional.

## Client configuration

Build the binary, then configure an MCP client to spawn it:

```json
{
  "mcpServers": {
    "bifunctional-org": {
      "command": "/absolute/path/to/repo/mcp/bin/bifunctional-mcp",
      "args": ["--config", "/absolute/path/to/repo/mcp/config/repos.local.json"]
    }
  }
}
```

The MCP process inherits `gh` authentication from the client process. Do not place tokens in this configuration.

## Adding a repository

1. Clone the repository locally.
2. Add it to the private local registry first.
3. Set the smallest useful `docsRoots` list.
4. Begin with `repo:read`, `docs:read`, and `issues:read` only.
5. Add write permissions only after read behavior is verified.
6. Run `make check` and call `list_repos`, `get_repo_context`, and `search_docs`.

## Failure behavior

- Invalid or unknown registry fields stop startup.
- Missing repositories stop startup instead of silently weakening policy.
- Missing docs roots are skipped during search, allowing staged wiki adoption.
- Oversized documents are rejected.
- Failed external writes are audited without recording issue bodies or secrets.
- Existing scaffold destinations are never overwritten.

## Recovery

- A failed scaffold is assembled in a temporary sibling directory and removed before return.
- A successful scaffold is a normal git repository with an initial commit and can be deleted or moved normally.
- A failed GitHub publication may leave the local scaffold intact. Inspect `git remote -v` and the GitHub organization before retrying.
- Rotate credentials immediately if they appear in git history, terminal output, or an issue.
