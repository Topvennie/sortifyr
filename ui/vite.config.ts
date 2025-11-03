import tailwindcss from "@tailwindcss/vite";
import react from '@vitejs/plugin-react-swc';
import path from 'path';
import { defineConfig } from 'vite';
import { imagetools } from 'vite-imagetools';
import checker from 'vite-plugin-checker';

export default defineConfig(() => {
  return {
    build: {
      outDir: '../public',
      emptyOutDir: true,
      sourcemap: true,
    },
    plugins: [
      tailwindcss(),
      react(),
      checker({
        // e.g. use TypeScript check
        typescript: true,
      }),
      imagetools(),
    ],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    server: {
      port: 3000,
      proxy: {
        '/api': {
          target: 'http://backend:3001',
          changeOrigin: true,
        }
      }
    }
  }
})
