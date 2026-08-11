# Agent instructions

This repository has two deliberately separate products:

- `src/` and `public/`: the public Bifunctional Astro website.
- `mcp/`: the local-first organization knowledge and repository control plane.
- `skills/diagram-design/`, `.codex-plugin/`, `commands/`, and `diagrams/`: the project-local editorial diagram system and its branded working area.

## Before editing

- Read `CONTRIBUTING.md` for repository-wide practices.
- For website work, read the brand and theme files named there.
- For MCP work, read `mcp/README.md` and `colabspace/org/security-policy.md`.
- Treat `mcp/config/repos.json` as the effective access-control policy.

## Safety

- Never add credentials, client data, private exports, or personal information.
- Do not broaden repository roots or write permissions as a side effect of another change.
- Keep public-site content independent from `colabspace/` and `mcp/`.
- Keep diagram outputs free of private source material unless they are explicitly scoped for a private documentation root.
- Preserve explicit confirmation gates for filesystem and GitHub writes.
- Add tests for path handling, authorization, redaction, and external command changes.

## Verification

Website changes:

```bash
npm run build
```

MCP changes:

```bash
cd mcp
make check
```

Run both when changing shared repository configuration or CI.
