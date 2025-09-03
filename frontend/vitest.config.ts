import { defineConfig } from 'vitest/config';
import path from 'node:path';

export default defineConfig({
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/setupTests.ts'],
    globals: true,
    css: false,
  },
  resolve: {
    alias: {
      '@src': path.resolve(__dirname, 'src'),
      '@': path.resolve(__dirname, './src'),
    },
  },
});
