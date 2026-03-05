import { defineConfig } from 'astro/config';

export default defineConfig({
  site: 'https://bifunctional.xyz',
  base: '/',
  output: 'static',
  build: {
    format: 'file', // Output blog.html (not blog/index.html) — preserves original URL structure
  },
});
