# Bifunctional diagrams

This repository includes the MIT-licensed [diagram-design system](https://github.com/cathrynlavery/diagram-design) as a project-local Codex plugin. It provides editorial architecture, flowchart, sequence, data, planning, and systems diagrams as standalone HTML with inline SVG.

The upstream grammar is inherited; the visual skin is Bifunctional-owned and lives in [`skills/diagram-design/references/style-guide.md`](../skills/diagram-design/references/style-guide.md). It maps the website tokens to semantic diagram roles:

- cream paper with evergreen ink
- emerald structure and muted labels
- bubblegum pink as the single focal accent
- Fraunces for titles and callouts
- Montserrat for node names and labels
- system mono only for technical sublabels

## Working convention

Keep generated or hand-authored diagram HTML in this directory unless it is intentionally being published as website content. Each diagram should stay self-contained, use inline SVG, and follow the system's one-accent, 4px-grid, no-shadow, and orthogonal-connector rules.

Use the project-local `diagram-design` skill for new diagrams. Export a finished HTML diagram with `commands/export-diagram.md` when an SVG or PNG is needed.

The vendored upstream files and third-party attribution are kept under [`skills/diagram-design/`](../skills/diagram-design/). The local integration is intentionally limited to the plugin manifest, the Bifunctional skin, and the branded template/gallery assets.

The current upstream snapshot is commit `c5805a0` from the [diagram-design repository](https://github.com/cathrynlavery/diagram-design).
