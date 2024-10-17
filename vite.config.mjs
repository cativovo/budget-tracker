import { defineConfig } from "vite";
import tailwindcss from "tailwindcss";

// vite.config.js
export default defineConfig({
  css: {
    postcss: {
      plugins: [tailwindcss()],
    },
  },
  build: {
    manifest: true,
    rollupOptions: {
      input: "js/index.js",
    },
  },
});
