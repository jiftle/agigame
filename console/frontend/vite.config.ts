import react from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig, loadEnv } from 'vite';
import { fileURLToPath, URL } from 'node:url';

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');

  return {
    plugins: [react(), tailwindcss()],
    optimizeDeps: {
      include: [
        '@ant-design/plots',
        '@ant-design/pro-components',
        '@ant-design/icons',
        'dayjs',
      ],
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    server: {
      port: Number(process.env.VITE_PORT || env.VITE_PORT) || 8001,
      proxy: {
        '/api': {
          target:
            process.env.VITE_PROXY_TARGET ||
            env.VITE_PROXY_TARGET ||
            'http://127.0.0.1:8000',
          changeOrigin: true,
        },
        '/ws': {
          target:
            process.env.VITE_PROXY_TARGET ||
            env.VITE_PROXY_TARGET ||
            'http://127.0.0.1:8000',
          changeOrigin: true,
          ws: true,
        },
      },
    },
    build: {
      outDir: 'dist',
      sourcemap: false,
      chunkSizeWarningLimit: 1600,
      rollupOptions: {
        output: {
          manualChunks: {
            react: ['react', 'react-dom', 'react-router-dom'],
            antd: ['antd', '@ant-design/icons'],
            pro: ['@ant-design/pro-components'],
            charts: ['@ant-design/plots'],
          },
        },
      },
    },
  };
});
