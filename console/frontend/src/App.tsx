import { QueryClientProvider } from '@tanstack/react-query';
import { App as AntdApp, ConfigProvider, theme } from 'antd';
import { useEffect, useSyncExternalStore, type ReactNode } from 'react';
import { BrowserRouter, useNavigate } from 'react-router-dom';

import { queryClient } from '@/api/queryClient';
import ErrorBoundary from '@/components/ErrorBoundary';
import AppRoutes from '@/router';
import { getSettings, subscribeSettings } from '@/theme';
import { setAntdFeedback } from '@/utils/antdApp';
import { setNavigator } from '@/utils/history';

/** 把 react-router 的 navigate 注入到薄封装，供非组件模块使用 */
function HistoryBridge() {
  const navigate = useNavigate();
  useEffect(() => {
    setNavigator((to, options) => navigate(to, options));
  }, [navigate]);
  return null;
}

/** 把 antd App 上下文（含主题）暴露给非组件模块（如 axios 拦截器） */
function AntdAppBridge() {
  const app = AntdApp.useApp();
  useEffect(() => {
    setAntdFeedback(app);
  }, [app]);
  return null;
}

function ThemeProvider({ children }: { children: ReactNode }) {
  const settings = useSyncExternalStore(subscribeSettings, getSettings);
  const isDark = settings.navTheme === 'realDark';
  return (
    <ConfigProvider
      variant="filled"
      theme={{
        algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm,
        token: {
          colorPrimary: settings.colorPrimary || '#1677ff',
          fontFamily: 'AlibabaSans, sans-serif',
        },
      }}
    >
      <AntdApp>
        <AntdAppBridge />
        {children}
      </AntdApp>
    </ConfigProvider>
  );
}

export default function App() {
  return (
    <ErrorBoundary>
      <QueryClientProvider client={queryClient}>
        <ThemeProvider>
          <BrowserRouter>
            <HistoryBridge />
            <AppRoutes />
          </BrowserRouter>
        </ThemeProvider>
      </QueryClientProvider>
    </ErrorBoundary>
  );
}
