import {
  ApiOutlined,
  FullscreenExitOutlined,
  FullscreenOutlined,
  MoonOutlined,
  ReloadOutlined,
  SunOutlined,
} from '@ant-design/icons';
import { ProLayout, SettingDrawer } from '@ant-design/pro-components';
import { Button, theme } from 'antd';
import { Suspense, useEffect, useMemo, useState, useSyncExternalStore } from 'react';
import { Link, Outlet, useLocation } from 'react-router-dom';

import { AccessProvider, buildAccess, type AccessMap } from '@/access';
import AvatarDropdown from '@/components/AvatarDropdown';
import { menuConfig, type MenuItemConfig } from '@/config/menu';
import { useCurrentUser } from '@/hooks/useCurrentUser';
import GlobalLoading from '@/loading';
import {
  getSettings,
  setSettings,
  subscribeSettings,
  type LayoutSettings,
} from '@/theme';

const isDev = import.meta.env.DEV;

function filterByAccess(items: MenuItemConfig[], access: AccessMap): MenuItemConfig[] {
  const result: MenuItemConfig[] = [];
  for (const item of items) {
    if (item.access && !access[item.access]) {
      continue;
    }
    if (item.children) {
      const children = filterByAccess(item.children, access);
      if (children.length === 0) {
        continue;
      }
      result.push({ ...item, children });
    } else {
      result.push(item);
    }
  }
  return result;
}

function toProRoutes(items: MenuItemConfig[]): Record<string, unknown>[] {
  return items.map((item) => ({
    path: item.path,
    name: item.name,
    icon: item.icon,
    routes: item.children ? toProRoutes(item.children) : undefined,
  }));
}

function toggleFullscreen() {
  if (document.fullscreenElement) {
    document.exitFullscreen();
  } else {
    document.documentElement.requestFullscreen();
  }
}

export default function BasicLayout() {
  const location = useLocation();
  const settings = useSyncExternalStore(subscribeSettings, getSettings);
  const { data } = useCurrentUser();
  const { token } = theme.useToken();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [refreshTick, setRefreshTick] = useState(0);
  const [isFullscreen, setIsFullscreen] = useState(false);

  const isDark = settings.navTheme === 'realDark';

  useEffect(() => {
    const onChange = () => setIsFullscreen(!!document.fullscreenElement);
    document.addEventListener('fullscreenchange', onChange);
    return () => document.removeEventListener('fullscreenchange', onChange);
  }, []);

  const toggleTheme = () =>
    setSettings({ ...settings, navTheme: isDark ? 'light' : 'realDark' });

  const access = useMemo(
    () => buildAccess(data?.perms),
    [data],
  );
  const menuRoutes = useMemo(
    () => toProRoutes(filterByAccess(menuConfig, access)),
    [access],
  );

  return (
    <AccessProvider access={access}>
      <ProLayout
        {...settings}
        title="AdminBase"
        logo="/logo.svg"
        className={`adminbase-layout ${
          isDark ? 'adminbase-layout--dark' : 'adminbase-layout--light'
        }`}
        route={{ routes: menuRoutes }}
        location={{ pathname: location.pathname }}
        itemRender={(route) =>
          route.path ? (
            <Link to={route.path}>{route.breadcrumbName}</Link>
          ) : (
            <span>{route.breadcrumbName}</span>
          )
        }
        menuItemRender={(item, dom) =>
          item.path ? <Link to={item.path}>{dom}</Link> : dom
        }
        actionsRender={() =>
          [
            isDev ? (
              <a
                key="api"
                href="http://127.0.0.1:8000/swagger"
                target="_blank"
                rel="noreferrer"
              >
                <ApiOutlined /> <span>API 文档</span>
              </a>
            ) : null,
            <Button
              key="refresh"
              type="text"
              icon={<ReloadOutlined />}
              title="刷新当前页"
              onClick={() => setRefreshTick((tick) => tick + 1)}
            />,
            <Button
              key="theme"
              type="text"
              icon={isDark ? <SunOutlined /> : <MoonOutlined />}
              title={isDark ? '切换到亮色模式' : '切换到暗色模式'}
              onClick={toggleTheme}
            />,
            <Button
              key="fullscreen"
              type="text"
              icon={isFullscreen ? <FullscreenExitOutlined /> : <FullscreenOutlined />}
              title="全屏"
              onClick={toggleFullscreen}
            />,
            <AvatarDropdown
              key="user"
              user={data?.user}
              onOpenTheme={() => setDrawerOpen(true)}
            />,
          ].filter(Boolean)
        }
        waterMarkProps={
          settings.waterMark
            ? {
                content: data?.user?.nickname || 'AdminBase',
                font: { color: token.colorTextQuaternary },
              }
            : undefined
        }
        footerRender={() => (
          <div
            style={{
              textAlign: 'center',
              color: token.colorTextTertiary,
              padding: '16px 0',
            }}
          >
            AdminBase © {new Date().getFullYear()} 飞鱼科技
          </div>
        )}
      >
        <div key={refreshTick}>
          <Suspense fallback={<GlobalLoading />}>
            <Outlet />
          </Suspense>
        </div>
      </ProLayout>
      <SettingDrawer
        disableUrlParams
        enableDarkTheme
        collapse={drawerOpen}
        onCollapseChange={setDrawerOpen}
        settings={settings as never}
        onSettingChange={(next) => {
          setSettings({ ...settings, ...(next as Partial<LayoutSettings>) });
        }}
      />
    </AccessProvider>
  );
}
