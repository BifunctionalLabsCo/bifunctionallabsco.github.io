---
title: Organization MCP roadmap
owner: bifunctional
status: active
updated: 2026-08-07
tags: [mcp, roadmap, okf]
---

# Organization MCP roadmap

## Delivered foundation

- Versioned repository registry and operation policy.
- Local git metadata and OKF documentation search.
- Path containment, symlink defense, exclusions, size limits, and redaction.
- GitHub repository, issue, label, milestone, issue-write, and pull-request tools.
- Socratic intake and three project design documents.
- Atomic repository scaffold, initial git commit, and GitHub publication.
- JSONL audit trail, tests, CI, static build, and operator documentation.

## Next

- Register private repositories through a reviewed local-only registry.
- Add richer ranking and incremental indexes after measuring repository volume.
- Add decision-record and release-note generators.
- Add optional issue-to-doc summaries with conflict and staleness markers.
- Add a dedicated private control-plane repository if public path metadata becomes limiting.

## Later

- Remote authenticated transport only if a concrete multi-user need appears.
- Organization-wide policy distribution and schema migrations.
- Observability that reports latency and failures without recording content.
- Multi-repository release and dependency views.

## Non-goals

- Becoming a general shell agent.
- Indexing every source file.
- Replacing git or GitHub as the system of record.
- Storing secrets.
- Publishing internal knowledge through the website.
