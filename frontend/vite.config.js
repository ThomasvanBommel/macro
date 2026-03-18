import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    outDir: 'build',
    emptyOutDir: true,
  },
  server: process.env.DEVMODE === 'true' ? {
    proxy: {
      '/api': {
        target: 'http://backend:8080',
        changeOrigin: true,
      }
    }
  }: {}
})
