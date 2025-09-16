import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';
import tailwindcss from '@tailwindcss/vite';
import { visualizer } from 'rollup-plugin-visualizer';
import Inspect from 'vite-plugin-inspect';
import path from 'path';

// https://vite.dev/config/
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    sourcemap: false,
  },
  plugins: [react(), tailwindcss(), visualizer({ open: false }), Inspect()],
  resolve: {
    alias: {
      '@src': path.resolve(__dirname, 'src'),
      '@': path.resolve(__dirname, './src'),
    },
  },
});
