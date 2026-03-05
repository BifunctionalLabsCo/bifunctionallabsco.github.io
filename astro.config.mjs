import { defineConfig } from 'astro/config';

export default defineConfig({
  site: 'https://BifunctionalLabsCo.github.io',
  // Project-site mode: base = '/<repo-name>'
  // Change to '/' if using a custom domain (and update site accordingly)
  base: '/studio-website',
  output: 'static',
  build: {
    format: 'file', // Output blog.html (not blog/index.html) — preserves original URL structure
  },
});
