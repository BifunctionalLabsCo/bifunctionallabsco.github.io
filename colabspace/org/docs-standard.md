---
title: OKF documentation standard
owner: bifunctional
status: active
updated: 2026-08-07
tags: [docs, standard, okf]
---

# OKF documentation standard

OKF is Bifunctional's operational knowledge format: a small Markdown convention for knowledge that remains readable by people, git, static tools, and MCP clients.

## Principles

- Knowledge stays with the repository that owns it.
- Markdown is canonical; generated indexes are secondary.
- One document owns each durable fact.
- Links replace duplicated prose.
- Decisions record rationale and consequences, not only outcomes.
- Working notes graduate into canonical documents or expire.

## Frontmatter

Every canonical document uses:

```yaml
---
title: Human-readable title
owner: accountable-team-or-person
status: draft
updated: 2026-08-07
tags: [topic, okf]
---
```

Allowed status values should be defined per repository. Recommended values are `draft`, `active`, `deprecated`, and `archived`.

## Repository structure

```text
docs/
  index.md
  architecture.md
  roadmap.md
  decisions/
  design/
    brand-direction.md
    design-system.md
  research/
    project-brief.md
  issues/
  notes/
```

## Document rules

- `docs/index.md` links to canonical entry points.
- `architecture.md` describes boundaries, dependencies, data, deployment, and failure behavior.
- `decisions/` contains immutable decision records; superseding decisions link backward.
- `design/brand-direction.md` explains narrative, inherited traits, distinct traits, voice, and guardrails.
- `design/design-system.md` explains foundations, tokens, component states, accessibility, and responsive behavior.
- `research/project-brief.md` defines audience, problem, outcome, constraints, workflows, and non-goals.
- `issues/` stores durable summaries only when GitHub issue history is insufficient.
- `notes/` is non-canonical and should not become a permanent dumping ground.

## Search behavior

The MCP server searches title, tags, and body. It returns repository-relative source paths. Files outside configured roots, secret-like files, oversized files, and excluded paths are not searchable.
