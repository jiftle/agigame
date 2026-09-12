import {
  EyeInvisibleOutlined,
  EyeOutlined,
  LockOutlined,
  LogoutOutlined,
  SkinOutlined,
  UserOutlined,
} from '@ant-design/icons';
import { Avatar, Dropdown } from 'antd';
import { useSyncExternalStore } from 'react';

import { queryClient } from '@/api/queryClient';
import { logout } from '@/services/auth';
import { useAuthStore } from '@/stores/auth';
import { getSettings, setSettings, subscribeSettings } from '@/theme';
import { history } from '@/utils/history';

interface Props {
  user?: {
    avatar?: string;
    nickname?: string;
  };
  onOpenTheme?: () => void;
}

export default function AvatarDropdown({ user, onOpenTheme }: Props) {
  const settings = useSyncExternalStore(subscribeSettings, getSettings);
  const waterMark = !!settings.waterMark;

  const toggleWaterMark = () => {
    setSettings({ ...getSettings(), waterMark: !getSettings().waterMark });
  };

  const handleLogout = () => {
    const { refreshToken, clear } = useAuthStore.getState();
    const done = () => {
      clear();
      queryClient.clear();
      history.push('/login');
    };
    if (refreshToken) {
      logout(refreshToken)
        .catch(() => {})
        .finally(done);
    } else {
      done();
    }
  };

  return (
    <Dropdown
      menu={{
        items: [
          { key: 'profile', icon: <UserOutlined />, label: '个人中心' },
          { key: 'password', icon: <LockOutlined />, label: '修改密码' },
          { key: 'theme', icon: <SkinOutlined />, label: '主题设置' },
          {
            key: 'watermark',
            icon: waterMark ? <EyeInvisibleOutlined /> : <EyeOutlined />,
            label: waterMark ? '隐藏水印' : '显示水印',
          },
          { type: 'divider' },
          { key: 'logout', icon: <LogoutOutlined />, label: '退出登录' },
        ],
        onClick: ({ key }) => {
          if (key === 'logout') {
            handleLogout();
          } else if (key === 'theme') {
            onOpenTheme?.();
          } else if (key === 'watermark') {
            toggleWaterMark();
          } else if (key === 'profile') {
            history.push('/account');
          } else if (key === 'password') {
            history.push('/account?tab=security');
          }
        },
      }}
    >
      <span className="adminbase-user-trigger">
        <Avatar size={28} src={user?.avatar || undefined} icon={<UserOutlined />} />
        {user?.nickname ? (
          <span className="adminbase-user-trigger-name">{user.nickname}</span>
        ) : null}
      </span>
    </Dropdown>
  );
}
