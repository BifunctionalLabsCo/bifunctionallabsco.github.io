---
title: Repository registry
owner: bifunctional
status: active
updated: 2026-08-07
tags: [repositories, policy, mcp, okf]
---

# Repository registry

The machine-readable registry at `mcp/config/repos.json` is the MCP server's source of truth for repository discovery and authorization.

## Required fields

| Field | Meaning |
| --- | --- |
| `slug` | Exact `organization/repository` GitHub identity. |
| `localPath` | Local clone used for filesystem and git reads. |
| `visibility` | Public, private, or internal classification. |
| `defaultBranch` | Expected integration branch. |
| `docsRoots` | Repository-relative roots that may be searched. |
| `sensitivity` | Public, internal, or restricted handling class. |
| `allowedOperations` | Explicit read and write capabilities. |
| `excludedPaths` | Subtrees denied even when nested under a docs root. |

## Operation vocabulary

- `repo:read`
- `docs:read`
- `issues:read`
- `issues:write`
- `pulls:read`
- `pulls:write`
- `branches:write`

Unknown operations grant no capability. Write operations also require a per-call confirmation value.

## Adoption process

1. Register a local clone in `repos.local.json`.
2. Grant only read operations.
3. Verify metadata and documentation results.
4. Review exclusions and sensitivity.
5. Add write operations only for a demonstrated workflow.
6. Move public-safe policy into the shared registry only when path disclosure is acceptable.
