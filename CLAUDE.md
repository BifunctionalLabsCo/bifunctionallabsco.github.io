# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

---

> This file is written to be useful to any AI coding agent, not only Claude Code.
> It is the authoritative in-repo context document. Keep it up to date when the architecture or conventions change.

---

## Project

**bifunctional.xyz** -- website for Bifunctional Labs Co, a Berlin-based founder-led freelance studio. Built with Astro v5, deployed to GitHub Pages on every push to `main` via `.github/workflows/deploy.yml`.

Live site: <https://bifunctional.xyz>

---

## Commands

```bash
npm install          # install deps (Astro only, no other runtime deps)
npm run dev          # dev server at http://localhost:4321/
npm run build        # static output to dist/ -- run this before finishing any code change
npm run preview      # preview the built output locally
```

No linter, no test suite. `npm run build` is the only required gate.

---

## Architecture

### Static site, no framework components

Pure Astro static output (`output: 'static'`). No React, no Vue, no islands. All interactivity is vanilla JS or Alpine.js loaded from CDN in `BaseLayout.astro`. The build produces flat HTML files.

### Layout hierarchy

```
BaseLayout.astro          -- outer shell: <head>, fonts (Google Fonts CDN), Alpine.js, SiteFooter
  PostLayout.astro        -- inner shell for all travel post pages (CSS + nav/hero/author/more-section)
    expeditions/*.astro   -- individual expedition posts (props + content slots only)
    posts/*.astro         -- legacy post pages (redirect to expeditions/)
```

`BaseLayout` accepts a `head` slot for per-page `<style is:inline>` blocks.  
`PostLayout` accepts named slots: `nav-left`, `nav-right`, `more-cards`. Default slot is the article body.

### CSS approach

**No shared stylesheet, no utility framework.** Every page that does not use `PostLayout` carries its own `<style is:inline>` block in `<Fragment slot="head">`. This is intentional -- pages are self-contained. When editing a page's styles, change only that file.

`PostLayout` is the single exception: its `<style is:inline>` is inherited by all post pages, so changes there affect all five expedition posts at once.

Reference files (not imported, not compiled -- read-only source of truth):
- `src/theme/tokens.css` -- CSS variable definitions
- `src/theme/global.css` -- base style reference

### Design tokens (from `tokens.css`)

```
--teal-400: #40ffb5    accent / highlights
--teal-600: #00634d
--teal-800: #032b22    primary dark, text on pink
--cream:    #efede7    page background (expeditions index, cards)
--pink:     #ffafda    "More Expeditions" section background on all post pages
--black:    #1b1b1b    body text
```

Fonts loaded from Google Fonts CDN:
- **Fraunces** -- serif/italic headings and pull quotes
- **Cormorant Garamond** -- caps-set body copy (`.caps`)
- **Montserrat** -- sans-serif nav, labels, UI
- **Poppins** -- display/hero titles

### Routes

| URL | File |
|---|---|
| `/` | `src/pages/index.astro` |
| `/expeditions/` | `src/pages/expeditions/index.astro` |
| `/expeditions/andaman/` | `src/pages/expeditions/andaman.astro` |
| `/expeditions/bangkok/` | `src/pages/expeditions/bangkok.astro` |
| `/expeditions/innsbruck/` | `src/pages/expeditions/innsbruck.astro` |
| `/expeditions/salzburg/` | `src/pages/expeditions/salzburg.astro` |
| `/expeditions/verona/` | `src/pages/expeditions/verona.astro` |
| `/posts/travelblogpost/` | redirect 301 -> `/expeditions/bangkok/` |
| `/impressum/` | `src/pages/impressum/index.astro` |
| `/privacy-policy/` | `src/pages/privacy-policy/index.astro` |
| `/legal-notice/` | `src/pages/legal-notice/index.astro` |

### Image placeholder system

Post pages use `.media.ph` divs with a `data-label` attribute where real photos go:

```html
<div class="media ph" data-label="Description of the photo that goes here"></div>
```

CSS in `PostLayout` renders the label as centered muted text via `::after { content: attr(data-label) }`. When adding a real photo, replace the div with:

```html
<div class="media">
  <img src="/assets/images/filename.jpg" alt="Alt text" />
</div>
```

All images live in `public/assets/images/` and are served at `/assets/images/filename.jpg`. Images are copyright Zubin John -- not open source. See GitHub issue #2 for the planned Cloudflare R2 migration.

### Social links (canonical)

| Platform | URL |
|---|---|
| LinkedIn | `https://www.linkedin.com/in/zubin-john-953b4241/` |
| Bluesky | `https://bsky.app/profile/bifunctional.bsky.social` |
| Instagram | `https://www.instagram.com/zrjohn/` |
| Email (footer) | `admin@bifunctional.xyz` |
| Email (posts) | `zubin@bifunctional.xyz` |

---

## Creating a new expedition post

1. Create `src/pages/expeditions/yourslug.astro`
2. Use `PostLayout` -- pass props (see existing posts for reference), fill `default` slot with article sections, fill `more-cards` slot with three `<a class="more-card">` links
3. Add a card for it in `src/pages/expeditions/index.astro` (`.story-grid`)
4. Update the `more-cards` slot on related posts to cross-link

---

## Conventions

- **No em dashes** in copy. Use commas, colons, hyphens, or `&mdash;` only when genuinely needed.
- **Absolute asset paths** (`/assets/...`), never relative (`../assets/...`), so sub-pages resolve correctly.
- **Favicon**: `href="/assets/favicon.ico.png"` -- absolute path required in `BaseLayout`.
- **`npm run build` before every commit.** The build is the only quality gate.
- Brand context lives in `colabspace/Bifunctional Brand Guidelines.md` and `CONTRIBUTING.md`. Read them before making visual or copy changes.

## Commit format

```
type: short subject in lowercase

Optional body.

Co-Authored-By: Claude <noreply@anthropic.com>
```

Types: `feat`, `fix`, `refactor`, `docs`, `chore`. Subject under 72 chars. Co-author trailer should reflect whatever agent made the change.
