---
title: Security and information policy
owner: bifunctional
status: active
updated: 2026-08-07
tags: [security, privacy, policy, okf]
---

# Security and information policy

## Classification

### Public

Material intentionally safe for a public repository, public website, public issue, or model-facing result.

### Internal

Organization operating knowledge that may be queried only from approved local or private repositories. It must not be mirrored into the public website.

### Restricted

Client confidences, credentials, personal data, contracts, financial details, security findings, and other need-to-know material. Restricted content should not be model-indexed unless a separately reviewed workflow requires it.

## Rules

- Repository classification is explicit in the registry.
- Documentation roots are allowlists, never inferred from the whole repository.
- Private repositories belong in an ignored local registry unless publishing their identity is intentional.
- Credentials are never stored in Markdown, registry files, MCP configuration, prompts, issues, or audit logs.
- Public brand material may inform project work; private strategy must remain in private knowledge roots.
- External writes require the user to understand the exact repository and action before confirmation.
- Audit metadata records the operation and target, not sensitive content bodies.

## Public website

The Astro build reads `src/` and `public/`. It does not read `colabspace/` or `mcp/`. A future build integration must preserve this separation and receive a dedicated security review.

## Incident response

1. Stop the MCP client or revoke the affected registry operation.
2. Rotate exposed credentials.
3. Remove public exposure using GitHub's supported sensitive-data process.
4. Preserve a minimal incident timeline without copying secrets.
5. Add a regression test or policy check before restoring access.
