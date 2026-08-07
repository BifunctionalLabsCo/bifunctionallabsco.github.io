# bifunctional

Website for [Bifunctional Labs Co](https://bifunctional.xyz), a Berlin-based freelance studio. Built with [Astro](https://astro.build) and deployed to GitHub Pages.

This repository also hosts the separately bounded [Bifunctional organization MCP control plane](mcp/README.md). The MCP service is not included in the public Astro build and indexes only documentation roots explicitly allowlisted in its registry.

## Local development

```bash
npm install
npm run dev        # → http://localhost:4321/
npm run build      # Static output → dist/
npm run preview    # Preview the built output locally
```

## Deployment

Every push to **main** automatically builds and deploys to GitHub Pages via `.github/workflows/deploy.yml`. The live site is at [bifunctional.xyz](https://bifunctional.xyz).

## Project structure

```
src/
  layouts/
    BaseLayout.astro              - shared shell: <head>, fonts, Alpine.js CDN, footer
  pages/
    index.astro                   - homepage (hero, services, bios, expeditions teaser)
    expeditions/
      index.astro                 - expeditions listing page
      innsbruck.astro             - Into the Alps: Innsbruck (Oct 2021)
      salzburg.astro              - Into the Alps: Salzburg (Oct 2021)
      verona.astro                - Verona: Of Bridges & Brawls (Oct 2020)
    posts/
      travelblogpost.astro        - Bangkok Airwaves (Mar 2019)
    privacy-policy/index.astro
    legal-notice/index.astro
    impressum/index.astro
  components/
    SiteFooter.astro
  theme/
    tokens.css                    - design tokens (CSS variables, reference only)
    global.css                    - global style reference
public/
  assets/
    images/                       - portraits and travel photos
    favicon.ico.png
    hero.png, logo.png, etc.
.github/workflows/deploy.yml      - CI/CD pipeline
mcp/                              - Go MCP server and repository control plane
colabspace/org/                   - OKF organization standards and policy
```

## Images and assets

Image assets in `public/assets/` are copyright Zubin John and are not covered by any open source licence. The code is open; the photos and brand visuals are not.

Expedition post pages currently use labeled placeholder divs where images will go. See [issue #2](https://github.com/BifunctionalLabsCo/bifunctionallabsco.github.io/issues/2) for the planned migration to Cloudflare R2.

## Contributing

See `CONTRIBUTING.md` before making changes. It covers brand context, coding practices, copy style, and commit message conventions.
