import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

const target = process.env.API_PROXY_TARGET || 'http://localhost:8080'

export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',
    proxy: {
      '/api': {
        target,
        changeOrigin: true,
      },
      '/ws': {
        target: target.replace('http://', 'ws://'),
        ws: true,
      },
    },
  },
})