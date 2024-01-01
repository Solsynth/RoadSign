import { fileURLToPath, URL } from "node:url"

import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import unocss from "unocss/vite"

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    unocss()
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url))
    }
  },
  server: {
    proxy: {
      "/webhooks": "http://127.0.0.1:81",
      "/cgi": "http://127.0.0.1:81"
    }
  }
})
