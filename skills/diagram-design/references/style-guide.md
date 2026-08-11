# Style Guide

**The single source of truth for colors, typography, and tokens.** Every diagram draws from this — not from hex values inlined in other reference files. If you want to change the visual skin of Schematic, change this file.

This project uses the Bifunctional skin: warm cream paper, evergreen ink, emerald structure, and a restrained bubblegum focal accent. The diagram grammar remains inherited from the upstream system; the palette and typography below are the project-owned layer. Change this file when the Bifunctional diagram skin evolves so every new diagram inherits the change without touching type-specific logic.

To generate your own from a website URL, see [`onboarding.md`](onboarding.md).

---

## Tokens

### Semantic roles

Every token is referred to by **semantic role**, not by its hex value. Type references (`type-*.md`) and SKILL.md say `accent`, not a hard-coded color.

| Role | Purpose | Default (light) | Default (dark) |
|---|---|---|---|
| `paper` | Page background, default node fill | `#efede7` (cream) | `#032b22` (evergreen) |
| `paper-2` | Diagram container bg, secondary fill | `#f7f5f0` | `#10231d` |
| `ink` | Primary text, primary stroke | `#032b22` (evergreen) | `#efede7` (cream) |
| `muted` | Secondary text, default arrow stroke | `#00634d` (emerald) | `#d8f2ea` (atmospheric) |
| `soft` | Sublabels, boundary labels | `#4b6f65` | `#a9c4b9` |
| `rule` | Hairline borders | `rgba(3,43,34,0.14)` | `rgba(239,237,231,0.16)` |
| `rule-solid` | Stronger borders, baselines | `#9eb3aa` | `rgba(158,179,170,0.34)` |
| `accent` | Focal / 1–2 max per diagram | `#ff4f81` (bubblegum) | `#ff83a7` |
| `accent-tint` | Fill for accent-bordered boxes | `rgba(255,79,129,0.12)` | `rgba(255,131,167,0.16)` |
| `link` | HTTP/API calls, external arrows | `#a33c63` (deep blush) | `#ffafca` |

> **Bifunctional source:** this skin maps the project tokens from `src/theme/tokens.css` and `src/theme/global.css` — `cream #efede7`, `evergreen #032b22`, `emerald #00634d`, `tropical mint #40ffb5`, `blush #ffafda`, and `bubblegum #ff4f81`. The `soft`, `rule`, and `link` tokens are derived to preserve the diagram system's hierarchy and contrast rules.

> **Note:** The bundled templates, gallery, and example HTML files are checked against this Bifunctional skin. New diagrams the skill produces should use the semantic roles above rather than copying one-off hex values.

### Inversion rule (light → dark)

Any `rgba(3,43,34, X)` in light becomes `rgba(239,237,231, X)` in dark. Same opacities, RGB flipped. The accent gets a slight hue-shift brighter to read on dark paper.

### Series palette (multi-series chart types only)

A small set of desaturated, editorial-tone colors for chart types that genuinely need to distinguish multiple overlapping entities (currently: **radar**). The "1-focal" rule still holds — `accent` is reserved for the focal series; the palette below covers the rest.

| Token | Light | Dark | Notes |
|---|---|---|---|
| `series-1` | `#00634d` (emerald) | `#40a889` | Non-focal series |
| `series-2` | `#3e8f75` (muted mint) | `#78c4aa` | Non-focal series |
| `series-3` | `#c15b80` (muted pink) | `#e48baa` | Non-focal series |
| `series-4` | `#9c815b` (warm sand) | `#c8ac7f` | Non-focal series |
| `series-5` | `#5e6e69` (slate green) | `#91aaa1` | Non-focal series |

Fills sit at `0.18` opacity light, `0.22` dark; strokes use the full color. **Don't backfill these tokens to non-chart types** — architecture, swimlane, etc. continue to use muted-ink variants. The series palette is opt-in for diagrams where overlapping shapes demand distinguishable color, not a license to add color elsewhere.

### Terminal skin (opt-in alternate)

A self-contained palette for the terminal-window primitive (see [primitive-terminal.md](primitive-terminal.md)) — a CLI-chrome register for dev-tool posts and technical social cards. It does not replace the default skin above and isn't affected by onboarding; it's a second, fixed skin you opt into per-diagram.

| Token | Hex | Purpose |
|---|---|---|
| `terminal-page` | `#032b22` | Page background behind the window |
| `terminal-paper` | `#10231d` | Window body, node fill |
| `terminal-bar` | `#1a2f28` | Titlebar strip |
| `terminal-border` | `#36594d` | Window border, hairlines |
| `terminal-ink` | `#efede7` | Primary text, primary stroke |
| `terminal-muted` | `#9eb3aa` | Secondary text, sublabels, ring stroke |
| `terminal-soft` | `#5b766b` | Tertiary — inactive dots, spokes |
| `terminal-accent` | `#40ffb5` | The one accent — focal station, prompt sign, active dot |
| `terminal-accent-tint` | `rgba(64,255,181,0.12)` | Fill for accent-bordered boxes |

**1-accent rule still holds.** Everything that isn't `terminal-ink` or `terminal-muted`/`terminal-soft` should be `terminal-accent` — never introduce a second hue.

---

## Typography

| Role | Family | Size | Weight | Usage |
|---|---|---|---|---|
| `title` | Fraunces | 1.75rem | 400 | Page H1 |
| `node-name` | Montserrat (sans) | 12px | 600 | Human-readable labels |
| `sublabel` | System mono | 9px | 400 | Port, protocol, URL, field type |
| `eyebrow` | Montserrat | 7–8px | 600, tracked 0.18em, uppercase | Type tags, axis labels |
| `arrow-label` | System mono | 8px | 400, tracked 0.06em | Arrow annotations |
| `callout` | Fraunces *italic* | 14px | 400 | Editorial asides only |

### Font stack

```html
<link href="https://fonts.googleapis.com/css2?family=Fraunces:ital,opsz,wght@0,9..144,300;0,9..144,400;0,9..144,500;0,9..144,600;1,9..144,500&family=Montserrat:wght@400;500;600&display=swap" rel="stylesheet">
```

**Load-bearing rule:** Mono is for *technical* content (ports, commands, URLs, field types). Names go in Montserrat. Page title and editorial callouts use Fraunces. **Never use mono as a blanket "dev" font.**

---

## Stroke, radius, spacing

| Token | Value | Use |
|---|---|---|
| `stroke-thin` | `0.8` | Tag-box outlines, leaf nodes |
| `stroke-default` | `1` | Most strokes |
| `stroke-strong` | `1.2` | Emphasis strokes |
| `radius-sm` | `4` | Small tags |
| `radius-md` | `6` | Node boxes |
| `radius-lg` | `8` | Containers, rings |
| `grid` | `4` | Every coord, size, and gap is divisible by 4 (hard rule) |

---

## Node type → treatment

Semantic role combinations — reference these by name in type specs.

| Type | Fill | Stroke |
|---|---|---|
| `focal` (1–2 max) | `accent-tint` | `accent` |
| `backend` | `paper-2` | `ink` |
| `store` | `ink @ 0.05` | `muted` |
| `external` | `ink @ 0.03` | `ink @ 0.30` |
| `input` | `muted @ 0.10` | `soft` |
| `optional` | `ink @ 0.02` | `ink @ 0.20` dashed `4,3` |
| `security` | `accent @ 0.05` | `accent @ 0.50` dashed `4,4` |

---

## Customizing the skin

Three options:

1. **Run onboarding** — see [`onboarding.md`](onboarding.md). Drop a URL; the skill extracts the palette + fonts and rewrites this file.
2. **Edit by hand** — change the hex values in the tables above. Run the pre-output taste gate afterward to verify the accent still reads as "focal" against the new paper color.
3. **Brand handoff** — paste your existing design-token JSON into a new section here and map its tokens to the semantic roles above.

### Constraints (don't break these)

- **Contrast**: `ink` must hit WCAG AA on `paper`. `muted` must hit AA on `paper` for 11px+ text.
- **One accent**: `accent` is bubblegum pink. Tropical mint remains an opt-in terminal/secondary series color; do not use it as a second focal accent in standard diagrams.
- **No rainbow palette**: if your brand ships 8 colors, pick 3 (paper, ink, accent). The rest become `muted` variants.
- **Serif + sans + mono**: three roles, not more. Fraunces carries titles and callouts, Montserrat carries names and labels, and a system mono stack is reserved for technical sublabels.
- **Paper is warm-neutral, not pure white**: pure white turns the design sterile. Pick a cream, bone, or light grey with a hint of warmth.
- **Dot pattern is optional, not default**: the 22×22 dot pattern is an opt-in "dotted paper" variant (good for long-form editorial hero diagrams). The default background is a clean `paper` fill, no pattern. When the pattern is enabled, it should sit at ~10% opacity of `ink` on `paper` — visible but quiet.
- **Container is clean by default**: the diagram sits directly on the page paper, no secondary container background or border. A framed variant (`paper-2` bg + `rule` border + 8px radius + padding) is available as an opt-in for card-heavy layouts, but don't reach for it by default — the extra chrome fights the figure.
