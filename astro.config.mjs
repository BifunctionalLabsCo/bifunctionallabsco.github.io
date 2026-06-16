import { defineConfig } from 'astro/config';
import sitemap from '@astrojs/sitemap';

export default defineConfig({
  site: 'https://bifunctional.xyz',
  base: '/',
  output: 'static',
  integrations: [
    sitemap({
      filter: (page) => !page.includes('/posts/travelblogpost'),
    }),
  ],
});
