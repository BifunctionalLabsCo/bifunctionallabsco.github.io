# studio-website

A simple website for the Bifunctional Labs Co — built with [Astro](https://astro.build) and optionally enhanced with [Alpine.js](https://alpinejs.dev).

## Local development

```bash
npm install
npm run dev        # → http://localhost:4321/studio-website/
```

## Build & preview

```bash
npm run build      # Static output → dist/
npm run preview    # Preview the built site locally
```

## Deployment

Every push to **main** automatically builds and deploys to **GitHub Pages** via `.github/workflows/deploy.yml`.

### Setup (one-time)

1. Go to **Settings → Pages** in the GitHub repo.
2. Under **Build and deployment → Source**, select **GitHub Actions**.

### Custom domain (optional)

1. Add a `public/CNAME` file containing your domain (e.g. `www.example.com`).
2. In `astro.config.mjs`, change `site` to your domain and set `base: '/'`.
3. Configure DNS per [GitHub's docs](https://docs.github.com/en/pages/configuring-a-custom-domain-for-your-github-pages-site).

## Project structure

```
src/
  layouts/BaseLayout.astro   — shared <head>, fonts, Alpine.js CDN
  pages/index.astro          — homepage
  pages/blog.astro           — expeditions listing
  pages/posts/               — travel blog post templates
  theme/tokens.css           — design tokens (CSS variables)
  theme/global.css           — full homepage CSS (reference copy)
public/assets/               — images, icons, favicon, PDF
.github/workflows/deploy.yml — CI/CD pipeline
```

## Agentic Contribution And Coding Guidelines

Before applying any site changes, always read the local brand and code context first:

- Brand guidelines: `colabspace/Bifunctional Brand Guidelines.md`
- Design tokens: `src/theme/tokens.css`
- Global style reference: `src/theme/global.css`
- Shared page shell: `src/layouts/BaseLayout.astro`
- The exact page or component being changed, such as `src/pages/index.astro`

Use those files as the source of truth for palette, typography, spacing, tone, imagery, and existing Astro structure. Keep changes inside the current theme direction unless the task explicitly asks for a redesign.

## Alpine.js

Alpine.js is loaded globally via CDN. It is a no-op by default — interactive behaviors only activate when `x-data` attributes are present in the markup. The site renders identically with JS disabled.
