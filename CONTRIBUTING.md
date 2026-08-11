# Contribution Guidelines

This is a small Astro website for Bifunctional Labs. Keep changes focused, brand-aware, and easy to review.

## Before Editing

Always read the local context before changing the site:

- Brand guidelines: `colabspace/Bifunctional Brand Guidelines.md`
- Design tokens: `src/theme/tokens.css`
- Global style reference: `src/theme/global.css`
- Shared page shell: `src/layouts/BaseLayout.astro`
- The exact page or component being changed, such as `src/pages/index.astro`

Use these files as the source of truth for palette, typography, spacing, tone, imagery, and Astro structure. Stay inside the current theme direction unless a task explicitly asks for a redesign.

## Coding Practices

- Prefer the existing Astro, HTML, and CSS patterns over adding new abstractions.
- Keep edits scoped to the requested page, section, or component.
- Avoid new dependencies unless they clearly solve a real problem.
- Use existing CSS variables and brand colors before inventing new values.
- Keep responsive rules simple and check desktop, tablet-ish, and mobile widths when changing layout.
- Do not change unrelated content, routes, assets, or legal pages while doing visual work.

## Diagram work

The repository carries the MIT-licensed `diagram-design` system under `skills/diagram-design/`. Its layout grammar and type references are inherited from the upstream project; the project-owned visual layer is [`skills/diagram-design/references/style-guide.md`](skills/diagram-design/references/style-guide.md).

- Keep diagram source and generated exports in `diagrams/` unless they are intentionally published as website content.
- Change brand colors and typography in the style guide, not in individual type references.
- Use one focal accent, no shadows, a 4px grid, and orthogonal connectors.
- Do not include private client data, credentials, or internal control-plane details in public diagrams.
- Run `npm run build` before finishing code changes.

## Visual And Content Style

- Preserve the current editorial, founder-led, boutique studio feel.
- Avoid generic SaaS visuals, loud gradients, overdone effects, and decorative clutter.
- Make cards and operational content scannable. Left-align dense text by default.
- Keep CTAs clear and sparse.
- Use short, concrete copy. Avoid filler, inflated agency language, and vague innovation talk.
- Do not use em dashes unless they are grammatically necessary and genuinely improve the sentence. Prefer commas, periods, colons, or simple hyphens.
- Use Unicode only when the existing content or brand voice benefits from it. Otherwise keep text ASCII-friendly.

## Formatting

- Keep Markdown headings short and sentence-case or title-case.
- Use flat bullet lists. Avoid deep nesting.
- Use backticks for file paths, commands, CSS variables, and code names.
- Do not add long comments to code. Add a comment only when it prevents confusion.
- Keep line wrapping readable, especially in Markdown and long HTML copy.

## Commit Messages

Use conventional-style commit messages:

- `feat: widen offering cards`
- `fix: correct footer link spacing`
- `docs: add contribution guidelines`
- `chore: update build metadata`

Keep the subject short, lowercase after the prefix, and specific to the change. When Codex contributes, include this trailer in the commit body:

```text
Co-authored-by: Codex <codex@openai.com>
```

If a Codex-specific signature is requested, include it separately:

```text
Code-agent-signature: Codex
```
